package web

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/internal/handler"
	"github.com/zqz/web/backend/internal/handler/auth"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service"
)

// LayoutData is the data passed to the shared layout template.
type LayoutData struct {
	PageTitle            string
	Content              template.HTML
	TitleExtra           template.HTML // optional content to show to the right of the page title (e.g. files search)
	User                 *authUser
	ShowUsers            bool
	PublicUploadsEnabled bool
	MaxFileSizeMB        int64 // effective max upload size in MB for the current user; 0 = no limit (e.g. admin). Used on upload page.
}

type contextKey int

const contextKeyPublicUploads contextKey = 0

// PublicUploadsMiddleware injects the public_uploads_enabled setting into context so layout and pages can use it.
func PublicUploadsMiddleware(repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enabled := true
			if val, err := repo.Settings.Get(r.Context(), "public_uploads_enabled"); err == nil && val != "true" {
				enabled = false
			}
			ctx := context.WithValue(r.Context(), contextKeyPublicUploads, enabled)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// LayoutDataFromRequest builds layout data from the request context (auth + public uploads).
func LayoutDataFromRequest(r *http.Request) LayoutData {
	data := LayoutData{PublicUploadsEnabled: true}
	if v := r.Context().Value(contextKeyPublicUploads); v != nil {
		if b, ok := v.(bool); ok {
			data.PublicUploadsEnabled = b
		}
	}
	user := auth.GetUserFromContext(r.Context())
	if user != nil {
		textColour := "#0d0d0d"
		if user.Colour != "" {
			textColour = contrastTextColour(user.Colour)
		}
		data.User = &authUser{
			Name:       user.Name,
			Admin:      user.IsAdmin(),
			DisplayTag: user.DisplayTag,
			Colour:     user.Colour,
			TextColour: textColour,
			Banned:     user.Banned,
		}
		data.ShowUsers = user.IsAdmin()
	}
	return data
}

// PagesHandler serves layout-based pages (upload, file-edit, users, api-docs).
type PagesHandler struct {
	templates *template.Template
	userSvc   *service.UserService
	fileSvc   *service.FileService
}

// NewPagesHandler creates a handler that serves themed pages.
func NewPagesHandler(templates *template.Template, userSvc *service.UserService, fileSvc *service.FileService) *PagesHandler {
	return &PagesHandler{templates: templates, userSvc: userSvc, fileSvc: fileSvc}
}

// userFilesPageUser is the user display model for the user files page.
type userFilesPageUser struct {
	ID                   int32
	Name                 string
	Email                string
	Role                 string
	DisplayTag           string
	Colour               string
	TextColour           string // contrasting text colour for the tag
	Banned               bool
	MaxFileSizeOverrideMB int64 // 0 means use site default
}

// userFileRow is one file row for the user files page.
type userFileRow struct {
	Name        string
	Slug        string
	SizeFmt     string
	ContentType string
	Complete    bool
	DownloadURL string
	ViewURL     string
}

// Upload serves the upload (home) page.
func (h *PagesHandler) Upload(w http.ResponseWriter, r *http.Request) {
	data := LayoutDataFromRequest(r)
	data.PageTitle = ""
	var userID *int32
	isAdmin := false
	if user := auth.GetUserFromContext(r.Context()); user != nil {
		id := user.ID
		userID = &id
		isAdmin = user.IsAdmin()
	}
	if maxBytes, err := h.fileSvc.GetEffectiveMaxFileSize(r.Context(), userID, isAdmin); err == nil && maxBytes > 0 {
		data.MaxFileSizeMB = maxBytes / (1024 * 1024)
	}
	RenderLayoutWithData(w, h.templates, "content_upload", &data, r)
}

// Edit serves the file edit page (slug is read from URL in JS).
func (h *PagesHandler) Edit(w http.ResponseWriter, r *http.Request) {
	RenderLayout(w, h.templates, "content_file_edit", "page.file_edit", r)
}

// View serves the read-only file view page (slug from URL in JS). Access: guests = public, users = public + own, admins = all.
func (h *PagesHandler) View(w http.ResponseWriter, r *http.Request) {
	RenderLayout(w, h.templates, "content_file_view", "page.file_view", r)
}

// Users serves the users list page (admin only; middleware enforces).
func (h *PagesHandler) Users(w http.ResponseWriter, r *http.Request) {
	RenderLayout(w, h.templates, "content_users", "page.users", r)
}

// UserFiles serves the page for one user and their files (admin only). URL: /users/{id}.
func (h *PagesHandler) UserFiles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID := int32(id64)

	user, err := h.userSvc.GetUserByID(r.Context(), userID)
	if err != nil || user == nil {
		if acceptsHTML(r) {
			data := LayoutDataFromRequest(r)
			data.PageTitle = "page.not_found"
			var buf bytes.Buffer
			_ = h.templates.ExecuteTemplate(&buf, "content_user_files", struct {
				User  *userFilesPageUser
				Files []userFileRow
			}{nil, nil})
			data.Content = template.HTML(buf.String())
			handler.SetContentType(w, handler.ContentTypeHTML)
			w.WriteHeader(http.StatusNotFound)
			_ = h.templates.ExecuteTemplate(w, "layout.html", data)
			return
		}
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	files, _ := h.fileSvc.ListFilesByUserID(r.Context(), userID, 1000, 0)
	rows := make([]userFileRow, 0, len(files))
	for _, f := range files {
		viewURL := ""
		if strings.HasPrefix(f.ContentType, "image/") {
			viewURL = "/api/v1/files/" + f.Slug + "/view"
		}
		rows = append(rows, userFileRow{
			Name:        f.Name,
			Slug:        f.Slug,
			SizeFmt:     formatBytesForUserFiles(int64(f.Size)),
			ContentType: f.ContentType,
			Complete:    f.BytesReceived == f.Size,
			DownloadURL: "/api/v1/files/" + f.Slug,
			ViewURL:     viewURL,
		})
	}

	maxMB := int64(0)
	if user.MaxFileSizeOverride != nil && *user.MaxFileSizeOverride > 0 {
		maxMB = *user.MaxFileSizeOverride / (1024 * 1024)
		if maxMB == 0 {
			maxMB = 1
		}
	}
	pageUser := &userFilesPageUser{
		ID:                    user.ID,
		Name:                  user.Name,
		Email:                 user.Email,
		Role:                  user.Role,
		DisplayTag:            user.DisplayTag,
		Colour:                user.Colour,
		Banned:                user.Banned,
		MaxFileSizeOverrideMB: maxMB,
	}
	if pageUser.Colour == "" {
		pageUser.Colour = "var(--border)"
		pageUser.TextColour = "#0d0d0d"
	} else {
		pageUser.TextColour = contrastTextColour(pageUser.Colour)
	}

	data := LayoutDataFromRequest(r)
	data.PageTitle = strings.ToLower(user.Name)
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "content_user_files", struct {
		User           *userFilesPageUser
		Files          []userFileRow
		ShowBanOption  bool
	}{pageUser, rows, true}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Content = template.HTML(buf.String())
	handler.SetContentType(w, handler.ContentTypeHTML)
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

func formatBytesForUserFiles(n int64) string {
	if n == 0 {
		return "0 B"
	}
	const k = 1024
	units := []string{"B", "KB", "MB", "GB"}
	i := 0
	for n >= k && i < len(units)-1 {
		n /= k
		i++
	}
	if i == 0 {
		return strconv.FormatInt(n, 10) + " " + units[i]
	}
	return strconv.FormatFloat(float64(n), 'f', 2, 64) + " " + units[i]
}

// UserSetBan handles POST /users/{id}/ban and POST /users/{id}/unban (admin only). Redirects back to user page.
func (h *PagesHandler) UserSetBan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := auth.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID := int32(id64)
	banned := !strings.HasSuffix(r.URL.Path, "/unban")
	_, err = h.userSvc.SetBanned(r.Context(), userID, banned)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users/"+idStr, http.StatusSeeOther)
}

// UserSetMaxFileSize handles POST /users/{id}/max-file-size (admin only). Form: max_file_size_mb (empty = use default).
func (h *PagesHandler) UserSetMaxFileSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := auth.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID := int32(id64)
	var maxBytes *int64
	if mbStr := r.FormValue("max_file_size_mb"); mbStr != "" {
		if mb, err := strconv.ParseInt(mbStr, 10, 64); err == nil && mb > 0 {
			b := mb * 1024 * 1024
			maxBytes = &b
		}
	}
	_, err = h.userSvc.SetMaxFileSize(r.Context(), userID, maxBytes)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/users/"+idStr, http.StatusSeeOther)
}

// UserSetProfile handles POST /users/{id}/profile (admin only). Form: display_tag, colour (hex #RRGGBB). Updates the user's display tag and colour.
func (h *PagesHandler) UserSetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	user := auth.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	idStr := chi.URLParam(r, "id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	userID := int32(id64)
	tag := strings.TrimSpace(r.FormValue("display_tag"))
	colour := strings.TrimSpace(r.FormValue("colour"))
	if colour != "" && !strings.HasPrefix(colour, "#") {
		colour = "#" + colour
	}
	_, err = h.userSvc.UpdateProfile(r.Context(), userID, tag, colour)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/users/"+idStr, http.StatusSeeOther)
}

// APIDocs serves the API documentation page.
func (h *PagesHandler) APIDocs(w http.ResponseWriter, r *http.Request) {
	RenderLayout(w, h.templates, "content_api_docs", "page.api_docs", r)
}

// Profile serves the user profile page (display tag and colour). Requires auth.
func (h *PagesHandler) Profile(w http.ResponseWriter, r *http.Request) {
	RenderLayout(w, h.templates, "content_profile", "page.profile", r)
}

// NotFound serves the 404 page. Use for r.NotFound.
func (h *PagesHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	if !acceptsHTML(r) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	data := LayoutDataFromRequest(r)
	data.PageTitle = "not found"
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "content_notfound", data); err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	data.Content = template.HTML(buf.String())
	handler.SetContentType(w, handler.ContentTypeHTML)
	w.WriteHeader(http.StatusNotFound)
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

// Forbidden serves the 403 page. Use when the user has no access (e.g. RequireAdmin).
func (h *PagesHandler) Forbidden(w http.ResponseWriter, r *http.Request) {
	if !acceptsHTML(r) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	data := LayoutDataFromRequest(r)
	data.PageTitle = "page.forbidden"
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "content_forbidden", data); err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	data.Content = template.HTML(buf.String())
	handler.SetContentType(w, handler.ContentTypeHTML)
	w.WriteHeader(http.StatusForbidden)
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

// Unauthorized serves the 401 page. Use when the user is not logged in (e.g. RequireAuth).
func (h *PagesHandler) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if !acceptsHTML(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data := LayoutDataFromRequest(r)
	data.PageTitle = "page.unauthorized"
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "content_unauthorized", data); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	data.Content = template.HTML(buf.String())
	handler.SetContentType(w, handler.ContentTypeHTML)
	w.WriteHeader(http.StatusUnauthorized)
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

// acceptsHTML returns true if the request Accept header prefers text/html.
func acceptsHTML(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "text/html")
}

// RenderLayout renders a layout page with the given content block and title.
func RenderLayout(w http.ResponseWriter, templates *template.Template, contentName, pageTitle string, r *http.Request) {
	data := LayoutDataFromRequest(r)
	data.PageTitle = pageTitle
	RenderLayoutWithData(w, templates, contentName, &data, r)
}

// RenderLayoutWithData renders a layout page with the given content block and pre-filled layout data.
func RenderLayoutWithData(w http.ResponseWriter, templates *template.Template, contentName string, data *LayoutData, r *http.Request) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, contentName, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Content = template.HTML(buf.String())
	handler.SetContentType(w, handler.ContentTypeHTML)
	_ = templates.ExecuteTemplate(w, "layout.html", data)
}

package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/zqz/web/backend/internal/handler"
	"github.com/zqz/web/backend/internal/handler/auth"
	"github.com/zqz/web/backend/internal/repository"
)

const siteSettingPublicUploads = "public_uploads_enabled"
const siteSettingDefaultMaxFileSize = "default_max_file_size"
const siteSettingAPIRateLimitRPS = "api_rate_limit_rps"

// AdminHandler serves the admin panel (admin only).
type AdminHandler struct {
	repo      *repository.Repository
	templates *template.Template
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(repo *repository.Repository, templates *template.Template) *AdminHandler {
	return &AdminHandler{repo: repo, templates: templates}
}

// AdminPageData is the data for the admin panel.
type AdminPageData struct {
	LayoutData
	FileCount            int64
	TotalSize            int64
	TotalSizeFmt         string
	UserCount            int64
	BannedCount          int64
	PublicUploadsEnabled bool
	DefaultMaxFileSizeMB int64 // 0 means use fallback (100 MB)
	APIRateLimitRPS      int   // API rate limit (requests/sec); 0 = disabled. Default 10.
}

// Page serves GET /admin (admin panel). Caller should use RequireAdmin middleware or check admin in handler.
func (h *AdminHandler) Page(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	ctx := r.Context()

	fileCount, _ := h.repo.Files.Count(ctx)
	totalSize, _ := h.repo.Files.TotalSize(ctx)
	userCount, _ := h.repo.Users.Count(ctx)
	bannedCount, _ := h.repo.Users.CountBanned(ctx)

	publicUploads := true
	if val, err := h.repo.Settings.Get(ctx, siteSettingPublicUploads); err == nil && val != "true" {
		publicUploads = false
	}

	defaultMaxFileSizeMB := int64(100)
	if val, err := h.repo.Settings.Get(ctx, siteSettingDefaultMaxFileSize); err == nil && val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil && n > 0 {
			defaultMaxFileSizeMB = n / (1024 * 1024)
			if defaultMaxFileSizeMB == 0 && n > 0 {
				defaultMaxFileSizeMB = 1
			}
		}
	}

	apiRateLimitRPS := 10
	if val, err := h.repo.Settings.Get(ctx, siteSettingAPIRateLimitRPS); err == nil && val != "" {
		if n, err := strconv.Atoi(val); err == nil && n >= 0 {
			apiRateLimitRPS = n
		}
	}

	data := AdminPageData{
		LayoutData:           LayoutDataFromRequest(r),
		FileCount:            fileCount,
		TotalSize:            totalSize,
		TotalSizeFmt:         formatBytesForAdmin(totalSize),
		UserCount:            userCount,
		BannedCount:          bannedCount,
		PublicUploadsEnabled: publicUploads,
		DefaultMaxFileSizeMB: defaultMaxFileSizeMB,
		APIRateLimitRPS:      apiRateLimitRPS,
	}
	data.PageTitle = "page.admin"

	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "content_admin", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Content = template.HTML(buf.String())

	handler.SetContentType(w, handler.ContentTypeHTML)
	_ = h.templates.ExecuteTemplate(w, "layout.html", data)
}

// UpdateSettings handles POST /admin/settings (public uploads, default max file size).
func (h *AdminHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := "false"
	if r.FormValue("public_uploads") == "on" || r.FormValue("public_uploads") == "1" {
		value = "true"
	}
	if err := h.repo.Settings.Set(r.Context(), siteSettingPublicUploads, value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if mbStr := r.FormValue("default_max_file_size_mb"); mbStr != "" {
		if mb, err := strconv.ParseInt(mbStr, 10, 64); err == nil && mb > 0 {
			bytesVal := strconv.FormatInt(mb*1024*1024, 10)
			if err := h.repo.Settings.Set(r.Context(), siteSettingDefaultMaxFileSize, bytesVal); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if rpsStr := r.FormValue("api_rate_limit_rps"); rpsStr != "" {
		if rps, err := strconv.Atoi(rpsStr); err == nil && rps >= 0 {
			if err := h.repo.Settings.Set(r.Context(), siteSettingAPIRateLimitRPS, strconv.Itoa(rps)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func formatBytesForAdmin(n int64) string {
	if n == 0 {
		return "0 B"
	}
	const k = 1024
	units := []string{"B", "KB", "MB", "GB", "TB"}
	i := 0
	for n >= k && i < len(units)-1 {
		n /= k
		i++
	}
	if i == 0 {
		return strconv.FormatInt(n, 10) + " " + units[i]
	}
	return fmt.Sprintf("%.2f %s", float64(n), units[i])
}

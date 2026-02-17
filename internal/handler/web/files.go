package web

import (
	"bytes"
	"html/template"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/zqz/web/backend/internal/handler"
	"github.com/zqz/web/backend/internal/handler/auth"
	"github.com/zqz/web/backend/internal/service"
)

// contrastTextColour returns a readable text colour (#ffffff or "#0d0d0d") for the given hex background (#RRGGBB).
func contrastTextColour(hex string) string {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return "#0d0d0d"
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	lum := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
	if lum <= 0.45 {
		return "#ffffff"
	}
	return "#0d0d0d"
}

// FilesHandler serves the files page and list fragment for htmx.
type FilesHandler struct {
	fileSvc   *service.FileService
	templates *template.Template
}

// FileRow is the view model for one file row in the list.
type FileRow struct {
	Name        string
	Comment     string
	Slug        string
	Size        int64
	SizeFmt     string
	ContentType string
	Private     bool
	ViewURL     string
	DownloadURL string
	CanEdit     bool
	ShowDelete  bool
	Complete    bool
}

// NewFilesHandler creates a FilesHandler with parsed templates.
func NewFilesHandler(fileSvc *service.FileService, templates *template.Template) *FilesHandler {
	return &FilesHandler{fileSvc: fileSvc, templates: templates}
}

// Page serves the full files page (shell with htmx that loads the list).
func (h *FilesHandler) Page(w http.ResponseWriter, r *http.Request) {
	data := LayoutDataFromRequest(r)
	data.PageTitle = "page.files"
	initialQ := strings.TrimSpace(r.URL.Query().Get("q"))
	initialQEncoded := ""
	if initialQ != "" {
		initialQEncoded = url.QueryEscape(initialQ)
	}
	filesPageData := struct {
		LayoutData
		FilesSearch        string
		FilesSearchEncoded string
	}{LayoutData: data, FilesSearch: initialQ, FilesSearchEncoded: initialQEncoded}
	var searchBuf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&searchBuf, "partial_files_search", filesPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.TitleExtra = template.HTML(searchBuf.String())
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "content_files", filesPageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data.Content = template.HTML(buf.String())
	handler.SetContentType(w, handler.ContentTypeHTML)
	if err := h.templates.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

const (
	maxListLimit        = 2000
	defaultListLimit    = 50
	loadMoreListLimit   = 100
)

// parseListParams reads limit and offset from the request query.
func parseListParams(r *http.Request, defaultLimit int32) (limit, offset int32) {
	limit = defaultLimit
	if s := r.URL.Query().Get("limit"); s != "" {
		if l, err := strconv.ParseInt(s, 10, 32); err == nil && l > 0 {
			limit = int32(l)
		}
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}
	if s := r.URL.Query().Get("offset"); s != "" {
		if o, err := strconv.ParseInt(s, 10, 32); err == nil && o >= 0 {
			offset = int32(o)
		}
	}
	return limit, offset
}

// List returns the file list fragment (initial load or load-more with OOB).
func (h *FilesHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseListParams(r, defaultListLimit)
	search := strings.TrimSpace(r.URL.Query().Get("q"))

	userID := auth.GetUserIDFromContext(r.Context())
	user := auth.GetUserFromContext(r.Context())
	isAdmin := user != nil && user.IsAdmin()

	files, err := h.fileSvc.ListFiles(r.Context(), limit, offset, userID, isAdmin, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows := make([]FileRow, len(files))
	for i, f := range files {
		viewURL := ""
		if strings.HasPrefix(f.ContentType, "image/") {
			viewURL = "/api/v1/files/" + f.Slug + "/view"
		}
		canEdit := userID != nil && (isAdmin || (f.UserID != nil && *f.UserID == *userID))
		rows[i] = FileRow{
			Name:        f.Name,
			Comment:     f.Comment,
			Slug:        f.Slug,
			Size:        int64(f.Size),
			SizeFmt:     formatBytes(int64(f.Size)),
			ContentType: humanReadableContentType(f.ContentType),
			Private:     f.Private,
			ViewURL:     viewURL,
			DownloadURL: "/api/v1/files/" + f.Slug,
			CanEdit:     canEdit,
			ShowDelete:  isAdmin,
			Complete:    f.BytesReceived == f.Size,
		}
	}

	hasMore := len(files) == int(limit)
	nextOffset := offset + limit

	searchEncoded := ""
	if search != "" {
		searchEncoded = url.QueryEscape(search)
	}

	data := struct {
		Rows          []FileRow
		Offset        int32
		Limit         int32
		LoadMoreLimit int32
		HasMore       bool
		NextOffset    int32
		Search        string
		SearchEncoded string
	}{Rows: rows, Offset: offset, Limit: limit, LoadMoreLimit: loadMoreListLimit, HasMore: hasMore, NextOffset: nextOffset, Search: search, SearchEncoded: searchEncoded}

	handler.SetContentType(w, handler.ContentTypeHTML)
	if offset == 0 {
		if err := h.templates.ExecuteTemplate(w, "files_list.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Load more: return OOB fragments to append to list and replace load-more area
	if err := h.templates.ExecuteTemplate(w, "files_list_append.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type authUser struct {
	Name       string
	Admin      bool
	DisplayTag string
	Colour     string
	TextColour string // contrasting text colour when Colour is set
	Banned     bool
}

// humanReadableContentType converts MIME types into short labels (pdf, jpg, mp4, etc).
func humanReadableContentType(contentType string) string {
	ct := strings.TrimSpace(contentType)
	if ct == "" {
		return "unknown"
	}

	mediaType := strings.ToLower(ct)
	if parsed, _, err := mime.ParseMediaType(ct); err == nil && parsed != "" {
		mediaType = strings.ToLower(parsed)
	}

	switch mediaType {
	case "application/pdf":
		return "pdf"
	case "image/jpeg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/webp":
		return "webp"
	case "image/svg+xml":
		return "svg"
	case "video/mp4":
		return "mp4"
	case "video/webm":
		return "webm"
	case "audio/mpeg":
		return "mp3"
	case "application/zip":
		return "zip"
	case "application/x-7z-compressed":
		return "7z"
	case "application/json":
		return "json"
	case "text/plain":
		return "txt"
	case "text/csv":
		return "csv"
	case "application/msword":
		return "doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return "docx"
	case "application/vnd.ms-excel":
		return "xls"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return "xlsx"
	case "application/vnd.ms-powerpoint":
		return "ppt"
	case "application/vnd.openxmlformats-officedocument.presentationml.presentation":
		return "pptx"
	}

	if strings.HasPrefix(mediaType, "image/") || strings.HasPrefix(mediaType, "video/") || strings.HasPrefix(mediaType, "audio/") || strings.HasPrefix(mediaType, "text/") {
		if idx := strings.IndexByte(mediaType, '/'); idx >= 0 && idx+1 < len(mediaType) {
			return mediaType[idx+1:]
		}
	}
	return mediaType
}

func formatBytes(n int64) string {
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

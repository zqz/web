package v1

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/handler/auth"
	"github.com/zqz/web/backend/internal/service"
)

const maxListLimit = 2000

// FileHandler handles file-related HTTP requests
type FileHandler struct {
	fileSvc *service.FileService
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileSvc *service.FileService) *FileHandler {
	return &FileHandler{fileSvc: fileSvc}
}

// FileResponse represents a file in API responses
type FileResponse struct {
	ID            int32     `json:"id"`
	Name          string    `json:"name"`
	Hash          string    `json:"hash"`
	Slug          string    `json:"slug"`
	Size          int32     `json:"size"`
	ContentType   string    `json:"content_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	BytesReceived int32     `json:"bytes_received"`
	Private       bool      `json:"private"`
	Comment       string    `json:"comment,omitempty"`
	UserID        *int32    `json:"user_id,omitempty"`
	ViewURL       string    `json:"view_url,omitempty"`
	DownloadURL   string    `json:"download_url"`
}

// toFileResponse converts a domain file to API response
func toFileResponse(f *domain.File) FileResponse {
	resp := FileResponse{
		ID:            f.ID,
		Name:          f.Name,
		Hash:          f.Hash,
		Slug:          f.Slug,
		Size:          f.Size,
		ContentType:   f.ContentType,
		CreatedAt:     f.CreatedAt,
		UpdatedAt:     f.UpdatedAt,
		BytesReceived: f.BytesReceived,
		Private:       f.Private,
		Comment:       f.Comment,
		UserID:        f.UserID,
		DownloadURL:   "/api/v1/files/" + f.Slug,
	}

	// Only add view URL for images
	if strings.HasPrefix(f.ContentType, "image/") {
		resp.ViewURL = "/api/v1/files/" + f.Slug + "/view"
	}

	return resp
}

// sanitizeContentDispositionFilename removes or replaces characters that could
// cause HTTP header injection (CR, LF, double-quote, backslash, control chars).
// Returns a safe filename for use in Content-Disposition, max 255 chars.
func sanitizeContentDispositionFilename(name string) string {
	var b strings.Builder
	for _, r := range name {
		if r == '\r' || r == '\n' || r == '"' || r == '\\' || r == '/' || unicode.IsControl(r) {
			continue
		}
		b.WriteRune(r)
	}
	s := strings.TrimSpace(b.String())
	if len(s) > 255 {
		s = s[:255]
	}
	if s == "" {
		s = "download"
	}
	return s
}

// handleFileServiceError writes the appropriate HTTP error for known file service errors.
// Returns true if the error was handled (caller should return), false otherwise.
func handleFileServiceError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, service.ErrFileNotFound):
		Error(w, http.StatusNotFound, err)
	case errors.Is(err, service.ErrUnauthorized):
		Error(w, http.StatusForbidden, err)
	case errors.Is(err, service.ErrFileIncomplete):
		Error(w, http.StatusConflict, err)
	default:
		return false
	}
	return true
}

// handleCreateFileError writes the appropriate HTTP error for create/upload service errors.
// Returns true if the error was handled, false otherwise.
func handleCreateFileError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, service.ErrPublicUploadsDisabled):
		ErrorMessage(w, http.StatusForbidden, "public uploads are disabled")
	case errors.Is(err, service.ErrFileTooLarge):
		ErrorMessage(w, http.StatusRequestEntityTooLarge, "file exceeds maximum allowed size")
	case errors.Is(err, service.ErrInvalidHash):
		ErrorMessage(w, http.StatusBadRequest, "hash must be a 64-character SHA-256 hex string")
	case errors.Is(err, service.ErrNameTooLong), errors.Is(err, service.ErrContentTypeTooLong):
		ErrorMessage(w, http.StatusBadRequest, err.Error())
	default:
		return false
	}
	return true
}

// streamFileWithHeaders sets response headers and streams the file body. disposition is "inline" or "attachment".
func streamFileWithHeaders(w http.ResponseWriter, r *http.Request, reader io.Reader, file *domain.File, disposition string) {
	safeName := sanitizeContentDispositionFilename(file.Name)
	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(int(file.Size)))
	w.Header().Set("Content-Disposition", disposition+`; filename="`+safeName+`"`)
	w.Header().Set("ETag", file.Hash)
	w.Header().Set("Cache-Control", "no-cache")
	if match := r.Header.Get("If-None-Match"); match == file.Hash {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	io.Copy(w, reader)
}

// parseListParams reads limit and offset from request query. defaultLimit is used when limit is missing or invalid.
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

// CreateFileRequest represents a file creation request.
// Private, comment, user_id, and slug are never taken from the body; they come from settings and auth.
type CreateFileRequest struct {
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	Size        int32  `json:"size"`
	ContentType string `json:"content_type"`
}

// CreateFile creates file metadata
func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request) {
	var req CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())
	user := auth.GetUserFromContext(r.Context())
	isAdmin := user != nil && user.IsAdmin()

	maxFileSize, err := h.fileSvc.GetEffectiveMaxFileSize(r.Context(), userID, isAdmin)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}

	file, err := h.fileSvc.CreateFile(r.Context(), domain.CreateFileRequest{
		Name:        req.Name,
		Hash:        req.Hash,
		Size:        req.Size,
		ContentType: req.ContentType,
		UserID:      userID,
		Private:     false,
		Comment:     "",
	}, maxFileSize)
	if err != nil {
		if handleCreateFileError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusCreated, toFileResponse(file))
}

// UploadFileData uploads file data (supports chunked uploads). Stops when max size exceeded; verifies SHA-256 on completion.
func (h *FileHandler) UploadFileData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		ErrorMessage(w, http.StatusBadRequest, "hash parameter is required")
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())
	user := auth.GetUserFromContext(r.Context())
	isAdmin := user != nil && user.IsAdmin()

	maxFileSize, err := h.fileSvc.GetEffectiveMaxFileSize(r.Context(), userID, isAdmin)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}
	if maxFileSize > 0 {
		r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
	}

	file, err := h.fileSvc.UploadFileData(r.Context(), hash, r.Body, maxFileSize)
	if err != nil {
		if errors.Is(err, service.ErrFileNotFound) {
			Error(w, http.StatusNotFound, err)
			return
		}
		if errors.Is(err, service.ErrHashMismatch) {
			ErrorMessage(w, http.StatusBadRequest, "file hash verification failed")
			return
		}
		if handleCreateFileError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, toFileResponse(file))
}

// GetFile retrieves file metadata by hash
func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		ErrorMessage(w, http.StatusBadRequest, "hash parameter is required")
		return
	}
	file, err := h.fileSvc.GetFileByHash(r.Context(), hash)
	if err != nil {
		if handleFileServiceError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, toFileResponse(file))
}

// ViewFile serves file for viewing (inline, images only)
func (h *FileHandler) ViewFile(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		ErrorMessage(w, http.StatusBadRequest, "slug parameter is required")
		return
	}
	userID := auth.GetUserIDFromContext(r.Context())
	reader, file, err := h.fileSvc.DownloadFile(r.Context(), slug, userID)
	if err != nil {
		if handleFileServiceError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	defer reader.Close()
	if !strings.HasPrefix(file.ContentType, "image/") {
		ErrorMessage(w, http.StatusBadRequest, "only images can be viewed inline")
		return
	}
	streamFileWithHeaders(w, r, reader, file, "inline")
}

// DownloadFile downloads file data by slug (with attachment content-disposition)
func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		ErrorMessage(w, http.StatusBadRequest, "slug parameter is required")
		return
	}
	userID := auth.GetUserIDFromContext(r.Context())
	reader, file, err := h.fileSvc.DownloadFile(r.Context(), slug, userID)
	if err != nil {
		if handleFileServiceError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	defer reader.Close()
	streamFileWithHeaders(w, r, reader, file, "attachment")
}

// ListFiles lists files with pagination
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	limit, offset := parseListParams(r, 50)
	search := strings.TrimSpace(r.URL.Query().Get("q"))
	userID := auth.GetUserIDFromContext(r.Context())
	user := auth.GetUserFromContext(r.Context())
	isAdmin := user != nil && user.IsAdmin()

	files, err := h.fileSvc.ListFiles(r.Context(), limit, offset, userID, isAdmin, search)
	if err != nil {
		Error(w, http.StatusInternalServerError, err)
		return
	}
	response := make([]FileResponse, len(files))
	for i, f := range files {
		response[i] = toFileResponse(f)
	}
	JSON(w, http.StatusOK, response)
}

// DeleteFile deletes a file
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		ErrorMessage(w, http.StatusBadRequest, "slug parameter is required")
		return
	}
	userID := auth.GetUserIDFromContext(r.Context())
	user := auth.GetUserFromContext(r.Context())
	isAdmin := user != nil && user.IsAdmin()

	err := h.fileSvc.DeleteFile(r.Context(), slug, userID, isAdmin)
	if err != nil {
		if handleFileServiceError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetFileBySlug retrieves file metadata by slug
func (h *FileHandler) GetFileBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		ErrorMessage(w, http.StatusBadRequest, "slug parameter is required")
		return
	}
	userID := auth.GetUserIDFromContext(r.Context())
	file, err := h.fileSvc.GetFileBySlug(r.Context(), slug, userID)
	if err != nil {
		if handleFileServiceError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, toFileResponse(file))
}

// UpdateFileRequest represents a file update request
type UpdateFileRequest struct {
	Name    *string `json:"name"`
	Private *bool   `json:"private"`
	Comment *string `json:"comment"`
}

// UpdateFile updates file metadata
func (h *FileHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if slug == "" {
		ErrorMessage(w, http.StatusBadRequest, "slug parameter is required")
		return
	}
	var req UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, err)
		return
	}
	userID := auth.GetUserIDFromContext(r.Context())
	user := auth.GetUserFromContext(r.Context())
	isAdmin := user != nil && user.IsAdmin()

	file, err := h.fileSvc.UpdateFile(r.Context(), slug, service.UpdateFileRequest{
		Name:    req.Name,
		Private: req.Private,
		Comment: req.Comment,
	}, userID, isAdmin)
	if err != nil {
		if handleFileServiceError(w, err) {
			return
		}
		Error(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, toFileResponse(file))
}

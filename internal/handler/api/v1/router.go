package v1

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/zqz/web/backend/internal/handler/auth"
)

const nonUploadTimeout = 200 * time.Millisecond

// timeoutForNonUpload cancels the request context after 200ms for all endpoints
// except POST /meta/{hash} (file data upload), which may take longer.
func timeoutForNonUpload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/meta/") {
			next.ServeHTTP(w, r)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), nonUploadTimeout)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewRouter creates a new API v1 router
func NewRouter(fileHandler *FileHandler, userHandler *UserHandler, authHandler *auth.AuthHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(timeoutForNonUpload)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // TODO: Configure allowed origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// File endpoints
	r.Route("/files", func(r chi.Router) {
		r.Post("/", fileHandler.CreateFile)         // Create file metadata
		r.Get("/", fileHandler.ListFiles)           // List files
		r.Get("/{slug}/view", fileHandler.ViewFile) // View file (inline, images only)
		r.Get("/{slug}", fileHandler.DownloadFile)  // Download file (attachment)
		r.Put("/{slug}", fileHandler.UpdateFile)    // Update file metadata
		r.Delete("/{slug}", fileHandler.DeleteFile) // Delete file
	})

	// File metadata endpoints (for web interface)
	r.Route("/file-metadata", func(r chi.Router) {
		r.Get("/{slug}", fileHandler.GetFileBySlug) // Get file metadata by slug
	})

	// Metadata endpoints
	r.Route("/meta", func(r chi.Router) {
		r.Get("/{hash}", fileHandler.GetFile)         // Get file metadata by hash
		r.Post("/{hash}", fileHandler.UploadFileData) // Upload file data
	})

	// User endpoints (admin only)
	r.Route("/users", func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return authHandler.RequireAdmin(next, nil) // API: no HTML error page
		})
		r.Get("/", userHandler.ListUsers)               // List all users
		r.Get("/{id}/files", userHandler.ListUserFiles) // List files by user
	})

	return r
}

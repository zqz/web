package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zqz/web/backend/internal/config"
	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/handler/auth"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service"
	"github.com/zqz/web/backend/internal/service/storage"
	"github.com/zqz/web/backend/internal/tests"
)

const testHash64 = "0000000000000000000000000000000000000000000000000000000000000001"

func setupFileHandlerTest(t *testing.T, ctx context.Context) (*chi.Mux, *service.FileService, func()) {
	t.Helper()

	pg, cleanup := tests.SetupTestDB(t, ctx)
	repo := repository.NewRepository(pg.Pool)

	err := repo.Settings.Set(ctx, "public_uploads_enabled", "true")
	require.NoError(t, err)

	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	fileSvc := service.NewFileService(repo, stor)
	userSvc := service.NewUserService(repo)

	logger := zerolog.Nop()
	cfg := &config.Config{SessionSecret: "test-secret", Env: "development"}
	authHandler := auth.NewAuthHandler(userSvc, &logger, cfg)
	fileHandler := NewFileHandler(fileSvc)
	userHandler := NewUserHandler(userSvc, fileSvc)

	r := chi.NewRouter()
	r.Use(authHandler.AuthMiddleware)
	r.Mount("/api/v1", NewRouter(fileHandler, userHandler, authHandler))

	return r, fileSvc, cleanup
}

func TestFileHandler_CreateFile_ListFiles_GetBySlug_Delete(t *testing.T) {
	ctx := context.Background()
	router, _, cleanup := setupFileHandlerTest(t, ctx)
	defer cleanup()

	// Create file metadata
	createBody := map[string]interface{}{
		"name":         "api-test.txt",
		"hash":         testHash64,
		"size":         100,
		"content_type": "text/plain",
		"private":      false,
		"comment":      "API test",
	}
	bodyBytes, _ := json.Marshal(createBody)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	var createResp FileResponse
	err := json.NewDecoder(rec.Body).Decode(&createResp)
	require.NoError(t, err)
	assert.NotZero(t, createResp.ID)
	assert.Equal(t, "api-test.txt", createResp.Name)
	assert.Equal(t, testHash64, createResp.Hash)
	assert.NotEmpty(t, createResp.Slug)
	assert.Equal(t, int32(100), createResp.Size)

	// List files
	req = httptest.NewRequest(http.MethodGet, "/api/v1/files/", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	var listResp []FileResponse
	err = json.NewDecoder(rec.Body).Decode(&listResp)
	require.NoError(t, err)
	require.Len(t, listResp, 1)
	assert.Equal(t, createResp.ID, listResp[0].ID)

	// Get file metadata by slug (before upload complete, slug is the initial one)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/file-metadata/"+createResp.Slug, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	var metaResp FileResponse
	err = json.NewDecoder(rec.Body).Decode(&metaResp)
	require.NoError(t, err)
	assert.Equal(t, createResp.ID, metaResp.ID)

	// Delete file (no auth = no ownership, so we expect 403 unless we allow anonymous delete for unowned? Check: DeleteFile checks ownership; anonymous can't delete a file they don't own. So we get 403.
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/files/"+createResp.Slug, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)

	// List still has the file
	req = httptest.NewRequest(http.MethodGet, "/api/v1/files/", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	err = json.NewDecoder(rec.Body).Decode(&listResp)
	require.NoError(t, err)
	assert.Len(t, listResp, 1)
}

func TestFileHandler_CreateFile_PublicUploadsDisabled(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	// Disable public uploads (migration inserts 'true' by default)
	err := repo.Settings.Set(ctx, "public_uploads_enabled", "false")
	require.NoError(t, err)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	fileSvc := service.NewFileService(repo, stor)
	userSvc := service.NewUserService(repo)
	logger := zerolog.Nop()
	cfg := &config.Config{SessionSecret: "test-secret", Env: "development"}
	authHandler := auth.NewAuthHandler(userSvc, &logger, cfg)
	r := chi.NewRouter()
	r.Use(authHandler.AuthMiddleware)
	r.Mount("/api/v1", NewRouter(NewFileHandler(fileSvc), NewUserHandler(userSvc, fileSvc), authHandler))

	createBody := map[string]interface{}{
		"name":         "anon.txt",
		"hash":         testHash64,
		"size":         10,
		"content_type": "text/plain",
		"private":      false,
		"comment":      "",
	}
	bodyBytes, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	var errResp ErrorResponse
	err = json.NewDecoder(rec.Body).Decode(&errResp)
	require.NoError(t, err)
	assert.Contains(t, errResp.Error, "public uploads")
}

// Upload + list flow is covered by service tests (file_service_test.go).
// API upload can hit 413 when GetEffectiveMaxFileSize differs in test DB; list/create/delete are covered above.

func TestFileHandler_GetFile_ByHash(t *testing.T) {
	ctx := context.Background()
	router, fileSvc, cleanup := setupFileHandlerTest(t, ctx)
	defer cleanup()

	file, err := fileSvc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "meta.txt",
		Hash:        testHash64,
		Size:        50,
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/meta/"+testHash64, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	var metaResp FileResponse
	err = json.NewDecoder(rec.Body).Decode(&metaResp)
	require.NoError(t, err)
	assert.Equal(t, file.ID, metaResp.ID)
	assert.Equal(t, testHash64, metaResp.Hash)
}

func TestFileHandler_GetFile_NotFound(t *testing.T) {
	ctx := context.Background()
	router, _, cleanup := setupFileHandlerTest(t, ctx)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/meta/0000000000000000000000000000000000000000000000000000000000009999", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestFileHandler_CreateFile_InvalidHash(t *testing.T) {
	ctx := context.Background()
	router, _, cleanup := setupFileHandlerTest(t, ctx)
	defer cleanup()

	createBody := map[string]interface{}{
		"name":         "bad.txt",
		"hash":         "not-64-hex",
		"size":         10,
		"content_type": "text/plain",
	}
	bodyBytes, _ := json.Marshal(createBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/files/", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestFileHandler_CreateFile_ValidationLimits(t *testing.T) {
	ctx := context.Background()
	router, _, cleanup := setupFileHandlerTest(t, ctx)
	defer cleanup()

	t.Run("name over 250 chars rejected", func(t *testing.T) {
		name := strings.Repeat("a", 251)
		createBody := map[string]interface{}{
			"name":         name,
			"hash":         testHash64,
			"size":         10,
			"content_type": "text/plain",
		}
		bodyBytes, _ := json.Marshal(createBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/files/", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("content_type over 80 chars rejected", func(t *testing.T) {
		contentType := strings.Repeat("a", 81)
		createBody := map[string]interface{}{
			"name":         "a.txt",
			"hash":         testHash64,
			"size":         10,
			"content_type": contentType,
		}
		bodyBytes, _ := json.Marshal(createBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/files/", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

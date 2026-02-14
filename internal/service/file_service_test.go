package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service/storage"
	"github.com/zqz/web/backend/internal/tests"
)

// valid SHA-256 hex hashes (64 chars) for tests that don't verify content
const testHash1 = "0000000000000000000000000000000000000000000000000000000000000001"
const testHash2 = "0000000000000000000000000000000000000000000000000000000000000002"
const testHashSame = "0000000000000000000000000000000000000000000000000000000000000abc"
const testHashPublic = "0000000000000000000000000000000000000000000000000000000000000def"
const testHashPrivate = "0000000000000000000000000000000000000000000000000000000000000bad"
const testHashPrivate2 = "0000000000000000000000000000000000000000000000000000000000000bed"

func TestFileService_CreateFile(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Test create file (0 = no size limit)
	file, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "test.txt",
		Hash:        testHash1,
		Size:        100,
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "Test file",
	}, 0)

	require.NoError(t, err)
	assert.NotZero(t, file.ID)
	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, testHash1, file.Hash)
	assert.NotEmpty(t, file.Slug)
	assert.Equal(t, int32(100), file.Size)
	assert.Equal(t, "text/plain", file.ContentType)
}

func TestFileService_CreateFile_Duplicate(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create first file
	file1, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "test.txt",
		Hash:        testHashSame,
		Size:        100,
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Try to create duplicate
	file2, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "duplicate.txt",
		Hash:        testHashSame, // Same hash
		Size:        100,
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)

	// Should return existing file
	require.NoError(t, err)
	assert.Equal(t, file1.ID, file2.ID)
	assert.Equal(t, file1.Hash, file2.Hash)
}

func TestFileService_UploadFileData(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create file metadata
	content := []byte("hello world")
	sum := sha256.Sum256(content)
	hash := fmt.Sprintf("%x", sum[:])

	file, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "upload.txt",
		Hash:        hash,
		Size:        int32(len(content)),
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Upload data
	reader := bytes.NewReader(content)
	uploadedFile, err := svc.UploadFileData(ctx, hash, reader, 0)
	require.NoError(t, err)

	assert.True(t, uploadedFile.Finished())
	assert.Equal(t, int32(len(content)), uploadedFile.BytesReceived)
	assert.NotEqual(t, file.Slug, uploadedFile.Slug) // Slug should be updated
}

func TestFileService_UploadFileData_Chunked(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create file metadata
	content := []byte("hello world")
	sum := sha256.Sum256(content)
	hash := fmt.Sprintf("%x", sum[:])

	_, err = svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "chunked.txt",
		Hash:        hash,
		Size:        int32(len(content)),
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Upload in chunks
	chunk1 := bytes.NewReader(content[:5])
	file1, err := svc.UploadFileData(ctx, hash, chunk1, 0)
	require.NoError(t, err)
	assert.False(t, file1.Finished())
	assert.Equal(t, int32(5), file1.BytesReceived)

	chunk2 := bytes.NewReader(content[5:])
	file2, err := svc.UploadFileData(ctx, hash, chunk2, 0)
	require.NoError(t, err)
	assert.True(t, file2.Finished())
	assert.Equal(t, int32(len(content)), file2.BytesReceived)
}

func TestFileService_GetFileBySlug_Public(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create public file
	created, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "public.txt",
		Hash:        testHashPublic,
		Size:        100,
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Anyone can access public file
	file, err := svc.GetFileBySlug(ctx, created.Slug, nil)
	require.NoError(t, err)
	assert.Equal(t, created.ID, file.ID)
}

func TestFileService_GetFileBySlug_Private_Unauthorized(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create owner user
	owner, err := repo.Users.Create(ctx, repository.CreateUserParams{
		Name:       "Owner",
		Email:      "owner@example.com",
		Provider:   "google",
		ProviderID: "owner-123",
		Role:       "member",
	})
	require.NoError(t, err)
	ownerID := owner.ID

	// Private comes from settings
	err = repo.Settings.Set(ctx, "default_private_upload", "true")
	require.NoError(t, err)
	created, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "private.txt",
		Hash:        testHashPrivate,
		Size:        100,
		ContentType: "text/plain",
		UserID:      &ownerID,
		Private:     false, // ignored; from settings
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Unauthenticated user cannot access
	_, err = svc.GetFileBySlug(ctx, created.Slug, nil)
	assert.ErrorIs(t, err, ErrUnauthorized)

	// Different user cannot access
	otherUserID := int32(2)
	_, err = svc.GetFileBySlug(ctx, created.Slug, &otherUserID)
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestFileService_GetFileBySlug_Private_Authorized(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create owner user
	owner, err := repo.Users.Create(ctx, repository.CreateUserParams{
		Name:       "Owner",
		Email:      "owner@example.com",
		Provider:   "google",
		ProviderID: "owner-123",
		Role:       "member",
	})
	require.NoError(t, err)
	ownerID := owner.ID

	err = repo.Settings.Set(ctx, "default_private_upload", "true")
	require.NoError(t, err)
	created, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "private.txt",
		Hash:        testHashPrivate,
		Size:        100,
		ContentType: "text/plain",
		UserID:      &ownerID,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Owner can access
	file, err := svc.GetFileBySlug(ctx, created.Slug, &ownerID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, file.ID)
}

func TestFileService_DownloadFile(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create and upload file
	content := []byte("download me")
	sum := sha256.Sum256(content)
	hash := fmt.Sprintf("%x", sum[:])

	created, err := svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "download.txt",
		Hash:        hash,
		Size:        int32(len(content)),
		ContentType: "text/plain",
		UserID:      nil,
		Private:     false,
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// Upload data
	uploaded, err := svc.UploadFileData(ctx, hash, bytes.NewReader(content), 0)
	require.NoError(t, err)

	// Download file (use updated slug from upload)
	reader, file, err := svc.DownloadFile(ctx, uploaded.Slug, nil)
	require.NoError(t, err)
	defer reader.Close()

	// Verify content
	downloaded, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, content, downloaded)
	assert.Equal(t, created.Hash, file.Hash)
}

func TestFileService_ListFiles_Public(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := repository.NewRepository(pg.Pool)
	stor, err := storage.NewDiskStorage(t.TempDir())
	require.NoError(t, err)

	svc := NewFileService(repo, stor)

	// Create public and private files (use valid 64-char hex hashes)
	hashes := []string{testHash1, testHash2, "0000000000000000000000000000000000000000000000000000000000000003"}
	for i := 0; i < 3; i++ {
		_, err := svc.CreateFile(ctx, domain.CreateFileRequest{
			Name:        fmt.Sprintf("public%d.txt", i),
			Hash:        hashes[i],
			Size:        100,
			ContentType: "text/plain",
			UserID:      nil,
			Private:     false,
			Comment:     "",
		}, 0)
		require.NoError(t, err)
	}

	// Create owner user
	owner, err := repo.Users.Create(ctx, repository.CreateUserParams{
		Name:       "Owner",
		Email:      "owner@example.com",
		Provider:   "google",
		ProviderID: "owner-123",
		Role:       "member",
	})
	require.NoError(t, err)
	ownerID := owner.ID

	// Private comes from settings; set so this file is created private
	err = repo.Settings.Set(ctx, "default_private_upload", "true")
	require.NoError(t, err)
	_, err = svc.CreateFile(ctx, domain.CreateFileRequest{
		Name:        "private.txt",
		Hash:        testHashPrivate2,
		Size:        100,
		ContentType: "text/plain",
		UserID:      &ownerID,
		Private:     false, // ignored; taken from settings
		Comment:     "",
	}, 0)
	require.NoError(t, err)

	// List as unauthenticated - should only see public files
	files, err := svc.ListFiles(ctx, 10, 0, nil, false, "")
	require.NoError(t, err)
	assert.Len(t, files, 3) // Only public files
}

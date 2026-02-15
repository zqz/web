package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zqz/web/backend/internal/tests"
)

const contentTypePlain = "text/plain"

func TestFileRepositoryCreate(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create a file
	file, err := repo.Files.Create(ctx, CreateFileParams{
		Size:        1024,
		Name:        "test.txt",
		Alias:       "test",
		Hash:        "abc123",
		Slug:        "test-slug",
		ContentType: contentTypePlain,
		UserID:      nil,
		Private:     false,
		Comment:     "Test file",
	})

	require.NoError(t, err)
	assert.NotZero(t, file.ID)
	assert.Equal(t, int32(1024), file.Size)
	assert.Equal(t, "test.txt", file.Name)
	assert.Equal(t, "abc123", file.Hash)
	assert.Equal(t, "test-slug", file.Slug)
	assert.Equal(t, contentTypePlain, file.ContentType)
	assert.False(t, file.Private)
	assert.Equal(t, "Test file", file.Comment)
}

func TestFileRepositoryGetByID(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create a file
	created, err := repo.Files.Create(ctx, CreateFileParams{
		Size:        2048,
		Name:        "get-by-id.txt",
		Alias:       "get",
		Hash:        "def456",
		Slug:        "get-slug",
		ContentType: contentTypePlain,
		UserID:      nil,
		Private:     false,
		Comment:     "",
	})
	require.NoError(t, err)

	// Get the file by ID
	file, err := repo.Files.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, file.ID)
	assert.Equal(t, created.Hash, file.Hash)
	assert.Equal(t, created.Name, file.Name)

	// Test not found
	_, err = repo.Files.GetByID(ctx, 99999)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestFileRepositoryGetBySlug(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create a file
	created, err := repo.Files.Create(ctx, CreateFileParams{
		Size:        3072,
		Name:        "slug-test.txt",
		Alias:       "slug",
		Hash:        "ghi789",
		Slug:        "unique-slug",
		ContentType: contentTypePlain,
		UserID:      nil,
		Private:     false,
		Comment:     "",
	})
	require.NoError(t, err)

	// Get the file by slug
	file, err := repo.Files.GetBySlug(ctx, "unique-slug")
	require.NoError(t, err)
	assert.Equal(t, created.ID, file.ID)
	assert.Equal(t, created.Slug, file.Slug)

	// Test not found
	_, err = repo.Files.GetBySlug(ctx, "nonexistent-slug")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestFileRepositoryGetByHash(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create a file
	created, err := repo.Files.Create(ctx, CreateFileParams{
		Size:        4096,
		Name:        "hash-test.txt",
		Alias:       "hash",
		Hash:        "unique-hash-123",
		Slug:        "hash-slug",
		ContentType: contentTypePlain,
		UserID:      nil,
		Private:     false,
		Comment:     "",
	})
	require.NoError(t, err)

	// Get the file by hash
	file, err := repo.Files.GetByHash(ctx, "unique-hash-123")
	require.NoError(t, err)
	assert.Equal(t, created.ID, file.ID)
	assert.Equal(t, created.Hash, file.Hash)

	// Test not found
	_, err = repo.Files.GetByHash(ctx, "nonexistent-hash")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestFileRepositoryList(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create multiple files
	for i := 0; i < 5; i++ {
		_, err := repo.Files.Create(ctx, CreateFileParams{
			Size:        int32(1024 * (i + 1)),
			Name:        "file" + string(rune('a'+i)) + ".txt",
			Alias:       "file" + string(rune('a'+i)),
			Hash:        "hash" + string(rune('a'+i)),
			Slug:        "slug" + string(rune('a'+i)),
			ContentType: contentTypePlain,
			UserID:      nil,
			Private:     false,
			Comment:     "",
		})
		require.NoError(t, err)
	}

	// List files
	files, err := repo.Files.List(ctx, 10, 0)
	require.NoError(t, err)
	assert.Len(t, files, 5)

	// Test pagination
	files, err = repo.Files.List(ctx, 2, 0)
	require.NoError(t, err)
	assert.Len(t, files, 2)

	files, err = repo.Files.List(ctx, 2, 2)
	require.NoError(t, err)
	assert.Len(t, files, 2)
}

func TestFileRepositoryUpdate(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create a file
	created, err := repo.Files.Create(ctx, CreateFileParams{
		Size:        5120,
		Name:        "update-test.txt",
		Alias:       "update",
		Hash:        "update-hash",
		Slug:        "update-slug",
		ContentType: contentTypePlain,
		UserID:      nil,
		Private:     false,
		Comment:     "Original comment",
	})
	require.NoError(t, err)

	// Update the file
	newName := "updated.txt"
	newComment := "Updated comment"
	updated, err := repo.Files.Update(ctx, UpdateFileParams{
		ID:      created.ID,
		Name:    &newName,
		Comment: &newComment,
	})
	require.NoError(t, err)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "updated.txt", updated.Name)
	assert.Equal(t, "Updated comment", updated.Comment)
	assert.Equal(t, created.Hash, updated.Hash) // Hash should remain unchanged
}

func TestFileRepositoryDelete(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Create a file
	created, err := repo.Files.Create(ctx, CreateFileParams{
		Size:        6144,
		Name:        "delete-test.txt",
		Alias:       "delete",
		Hash:        "delete-hash",
		Slug:        "delete-slug",
		ContentType: contentTypePlain,
		UserID:      nil,
		Private:     false,
		Comment:     "",
	})
	require.NoError(t, err)

	// Delete the file
	err = repo.Files.Delete(ctx, created.ID)
	require.NoError(t, err)

	// Verify it's deleted
	_, err = repo.Files.GetByID(ctx, created.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestFileRepositoryCount(t *testing.T) {
	ctx := context.Background()
	pg, cleanup := tests.SetupTestDB(t, ctx)
	defer cleanup()

	repo := NewRepository(pg.Pool)

	// Initial count should be 0
	count, err := repo.Files.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// Create files
	for i := 0; i < 3; i++ {
		_, err := repo.Files.Create(ctx, CreateFileParams{
			Size:        int32(1024 * (i + 1)),
			Name:        "count" + string(rune('a'+i)) + ".txt",
			Alias:       "count" + string(rune('a'+i)),
			Hash:        "counthash" + string(rune('a'+i)),
			Slug:        "countslug" + string(rune('a'+i)),
			ContentType: contentTypePlain,
			UserID:      nil,
			Private:     false,
			Comment:     "",
		})
		require.NoError(t, err)
	}

	// Count should be 3
	count, err = repo.Files.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

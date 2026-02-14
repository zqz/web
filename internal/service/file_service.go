package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service/storage"
)

var sha256HexRegex = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)

const (
	maxNameLen        = 250
	maxAliasLen       = 250
	maxContentTypeLen = 80
	sha256HexLen      = 64
)

var (
	// ErrFileNotFound is returned when a file is not found
	ErrFileNotFound = errors.New("file not found")

	// ErrUnauthorized is returned when a user doesn't have access to a file
	ErrUnauthorized = errors.New("unauthorized access to file")

	// ErrHashMismatch is returned when the file hash doesn't match the expected hash
	ErrHashMismatch = errors.New("file hash mismatch")

	// ErrFileIncomplete is returned when trying to access an incomplete file
	ErrFileIncomplete = errors.New("file upload incomplete")

	// ErrPublicUploadsDisabled is returned when anonymous uploads are disabled
	ErrPublicUploadsDisabled = errors.New("public uploads are disabled")

	// ErrFileTooLarge is returned when file size exceeds the allowed limit
	ErrFileTooLarge = errors.New("file exceeds maximum allowed size")

	// ErrInvalidHash is returned when the hash is not a valid SHA-256 hex string
	ErrInvalidHash = errors.New("hash must be a 64-character SHA-256 hex string")

	// ErrNameTooLong is returned when name exceeds max length
	ErrNameTooLong = errors.New("name must be at most 250 characters")

	// ErrContentTypeTooLong is returned when content type exceeds max length
	ErrContentTypeTooLong = errors.New("content type must be at most 80 characters")
)

// Processor defines the interface for file processing operations
type Processor interface {
	Process(ctx context.Context, file *domain.File, storage storage.Storage, repo *repository.Repository) error
	Name() string
}

// FileService handles file business logic
type FileService struct {
	repo       *repository.Repository
	storage    storage.Storage
	processors []Processor
}

// NewFileService creates a new file service
func NewFileService(repo *repository.Repository, storage storage.Storage) *FileService {
	return &FileService{
		repo:       repo,
		storage:    storage,
		processors: make([]Processor, 0),
	}
}

// AddProcessor adds a processor to run on file uploads
func (s *FileService) AddProcessor(p Processor) {
	s.processors = append(s.processors, p)
}

const settingDefaultMaxFileSize = "default_max_file_size"

// GetEffectiveMaxFileSize returns the max file size in bytes for the user. Returns 0 for no limit (admins).
func (s *FileService) GetEffectiveMaxFileSize(ctx context.Context, userID *int32, isAdmin bool) (int64, error) {
	if isAdmin {
		return 0, nil
	}
	var max int64
	if userID != nil {
		user, err := s.repo.Users.GetByID(ctx, *userID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return 0, fmt.Errorf("failed to get user: %w", err)
		}
		if err == nil && user.MaxFileSizeOverride != nil && *user.MaxFileSizeOverride > 0 {
			return *user.MaxFileSizeOverride, nil
		}
	}
	val, err := s.repo.Settings.Get(ctx, settingDefaultMaxFileSize)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return 0, fmt.Errorf("failed to get setting: %w", err)
	}
	if err == nil && val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil && n > 0 {
			max = n
		}
	}
	if max <= 0 {
		max = 100 * 1024 * 1024 // 100 MB fallback
	}
	return max, nil
}

// CreateFile creates a new file metadata entry. maxFileSize is the effective limit (0 = no limit).
func (s *FileService) CreateFile(ctx context.Context, req domain.CreateFileRequest, maxFileSize int64) (*domain.File, error) {
	// If anonymous upload, check if public uploads are enabled
	if req.UserID == nil {
		val, err := s.repo.Settings.Get(ctx, "public_uploads_enabled")
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, fmt.Errorf("failed to check setting: %w", err)
		}
		if err != nil || val != "true" {
			return nil, ErrPublicUploadsDisabled
		}
	}

	// Validate request (includes SHA-256 hash format)
	if err := validateCreateFileRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	if maxFileSize > 0 && int64(req.Size) > maxFileSize {
		return nil, ErrFileTooLarge
	}

	// Check if file with this hash already exists
	existing, err := s.repo.Files.GetByHash(ctx, req.Hash)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing file: %w", err)
	}

	if existing != nil {
		// File already exists, return it
		return dbFileToDoamin(existing), nil
	}

	// Private, comment, slug: from settings/server, not from user request
	private := false
	if v, err := s.repo.Settings.Get(ctx, "default_private_upload"); err == nil && v == "true" {
		private = true
	}
	comment := ""
	if v, err := s.repo.Settings.Get(ctx, "default_upload_comment"); err == nil {
		comment = v
	}
	alias := req.Name
	if len(alias) > maxAliasLen {
		alias = alias[:maxAliasLen]
	}
	slug := generateSlug(6)

	// Create file in repository
	dbFile, err := s.repo.Files.Create(ctx, repository.CreateFileParams{
		Size:          req.Size,
		Name:          req.Name,
		Alias:         alias,
		Hash:          req.Hash,
		Slug:          slug,
		ContentType:   req.ContentType,
		UserID:        req.UserID,
		Private:       private,
		Comment:       comment,
		BytesReceived: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return dbFileToDoamin(dbFile), nil
}

// UploadFileData uploads the actual file data. maxFileSize is the effective limit (0 = no limit).
// Stops reading as soon as max size would be exceeded and deletes partial data on overflow or hash mismatch.
func (s *FileService) UploadFileData(ctx context.Context, hash string, data io.Reader, maxFileSize int64) (*domain.File, error) {
	// Get file metadata
	dbFile, err := s.repo.Files.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	file := dbFileToDoamin(dbFile)

	// Check if already complete
	if file.Finished() {
		return file, nil
	}

	remainingToDeclared := int64(file.Size) - int64(file.BytesReceived)
	if remainingToDeclared <= 0 {
		return file, nil
	}

	var reader io.Reader = data
	if maxFileSize > 0 {
		allowed := maxFileSize - int64(file.BytesReceived)
		if allowed <= 0 {
			s.storage.Delete(hash)
			_, _ = s.repo.Files.Update(ctx, repository.UpdateFileParams{ID: dbFile.ID, BytesReceived: ptrInt32(0)})
			return nil, ErrFileTooLarge
		}
		limit := remainingToDeclared
		if allowed < limit {
			limit = allowed
		}
		reader = newMaxBytesReader(data, limit)
	}

	// Append data to storage
	bytesWritten, err := s.storage.Append(hash, reader)
	if err != nil {
		if errors.Is(err, ErrFileTooLarge) {
			s.storage.Delete(hash)
			_, _ = s.repo.Files.Update(ctx, repository.UpdateFileParams{ID: dbFile.ID, BytesReceived: ptrInt32(0)})
			return nil, ErrFileTooLarge
		}
		return nil, fmt.Errorf("failed to write file data: %w", err)
	}

	file.BytesReceived += int32(bytesWritten)

	// Update bytes_received in database
	bytesReceived := file.BytesReceived
	updatedFile, err := s.repo.Files.Update(ctx, repository.UpdateFileParams{
		ID:            dbFile.ID,
		BytesReceived: &bytesReceived,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update bytes_received: %w", err)
	}

	file = dbFileToDoamin(updatedFile)

	// If upload is complete, verify SHA-256 and delete on mismatch
	if file.Finished() {
		if err := s.verifyFileHash(hash); err != nil {
			s.storage.Delete(hash)
			return nil, err
		}

		// Update slug now that file is complete
		newSlug := generateSlug(6)
		updatedFile, err = s.repo.Files.Update(ctx, repository.UpdateFileParams{
			ID:   dbFile.ID,
			Slug: &newSlug,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update file slug: %w", err)
		}
		file = dbFileToDoamin(updatedFile)

		// Run processors
		if err := s.runProcessors(ctx, file); err != nil {
			// Log error but don't fail the upload
			// TODO: Add proper logging
			fmt.Printf("processor error: %v\n", err)
		}
	}

	return file, nil
}

// GetFileBySlug retrieves a file by its slug
func (s *FileService) GetFileBySlug(ctx context.Context, slug string, userID *int32) (*domain.File, error) {
	dbFile, err := s.repo.Files.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	file := dbFileToDoamin(dbFile)

	// Get current size from storage
	size, err := s.storage.Size(dbFile.Hash)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, fmt.Errorf("failed to get file size: %w", err)
	}
	if err == nil {
		file.BytesReceived = int32(size)
	}

	// Check access permissions
	if !file.CanBeAccessedBy(userID) {
		return nil, ErrUnauthorized
	}

	return file, nil
}

// GetFileByHash retrieves a file by its hash
func (s *FileService) GetFileByHash(ctx context.Context, hash string) (*domain.File, error) {
	dbFile, err := s.repo.Files.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	// Get current size from storage
	size, err := s.storage.Size(hash)
	if err != nil && !errors.Is(err, storage.ErrNotFound) {
		return nil, fmt.Errorf("failed to get file size: %w", err)
	}

	file := dbFileToDoamin(dbFile)
	if err == nil {
		file.BytesReceived = int32(size)
	}

	return file, nil
}

// DownloadFile returns a reader for downloading the file data
func (s *FileService) DownloadFile(ctx context.Context, slug string, userID *int32) (io.ReadCloser, *domain.File, error) {
	// Get file metadata and check permissions
	file, err := s.GetFileBySlug(ctx, slug, userID)
	if err != nil {
		return nil, nil, err
	}

	// Check if file is complete
	if !file.Finished() {
		return nil, nil, ErrFileIncomplete
	}

	// Get file data from storage
	reader, err := s.storage.Get(file.Hash)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, nil, ErrFileNotFound
		}
		return nil, nil, fmt.Errorf("failed to get file data: %w", err)
	}

	return reader, file, nil
}

// ListFiles returns a paginated list of files visible to the caller.
// If search is non-empty, filters by fuzzy match on name, alias, and comment (case-insensitive, via pg_trgm).
// Admins see all files; logged-in users see public files + their own; guests see only public.
func (s *FileService) ListFiles(ctx context.Context, limit, offset int32, userID *int32, isAdmin bool, search string) ([]*domain.File, error) {
	search = strings.TrimSpace(search)

	var dbFiles []*repository.File
	var err error

	if search != "" {
		if isAdmin {
			dbFiles, err = s.repo.Files.SearchFiles(ctx, search, limit, offset)
		} else if userID != nil {
			dbFiles, err = s.repo.Files.SearchFilesVisibleToUser(ctx, *userID, search, limit, offset)
		} else {
			dbFiles, err = s.repo.Files.SearchPublicFiles(ctx, search, limit, offset)
		}
	} else {
		if isAdmin {
			dbFiles, err = s.repo.Files.List(ctx, limit, offset)
		} else if userID != nil {
			dbFiles, err = s.repo.Files.ListVisibleToUser(ctx, *userID, limit, offset)
		} else {
			dbFiles, err = s.repo.Files.ListPublic(ctx, limit, offset)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	files := make([]*domain.File, len(dbFiles))
	for i, dbFile := range dbFiles {
		files[i] = dbFileToDoamin(dbFile)
	}

	return files, nil
}

// ListFilesByUserID returns a paginated list of files belonging to a specific user (for user profile / admin list).
func (s *FileService) ListFilesByUserID(ctx context.Context, userID int32, limit, offset int32) ([]*domain.File, error) {
	dbFiles, err := s.repo.Files.ListByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list files by user: %w", err)
	}
	files := make([]*domain.File, len(dbFiles))
	for i, dbFile := range dbFiles {
		files[i] = dbFileToDoamin(dbFile)
	}
	return files, nil
}

// UpdateFileRequest represents a request to update file metadata
type UpdateFileRequest struct {
	Name    *string
	Private *bool
	Comment *string
}

// UpdateFile updates file metadata. Owners and admins can update.
func (s *FileService) UpdateFile(ctx context.Context, slug string, req UpdateFileRequest, userID *int32, isAdmin bool) (*domain.File, error) {
	// Get file to check ownership
	file, err := s.GetFileBySlug(ctx, slug, userID)
	if err != nil {
		return nil, err
	}

	// Check if user owns the file or is admin
	if userID == nil || (!file.IsOwnedBy(*userID) && !isAdmin) {
		return nil, ErrUnauthorized
	}

	// Update file metadata
	updateParams := repository.UpdateFileParams{
		ID: file.ID,
	}

	if req.Name != nil {
		updateParams.Name = req.Name
	}
	if req.Private != nil {
		updateParams.Private = req.Private
	}
	if req.Comment != nil {
		updateParams.Comment = req.Comment
	}

	dbFile, err := s.repo.Files.Update(ctx, updateParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	return dbFileToDoamin(dbFile), nil
}

// DeleteFile deletes a file and its data
func (s *FileService) DeleteFile(ctx context.Context, slug string, userID *int32, isAdmin bool) error {
	// Get file to check ownership
	file, err := s.GetFileBySlug(ctx, slug, userID)
	if err != nil {
		return err
	}

	// Check if user owns the file or is admin
	if userID == nil || (!file.IsOwnedBy(*userID) && !isAdmin) {
		return ErrUnauthorized
	}

	// Delete thumbnail from database first (to avoid FK constraint violation)
	if err := s.repo.Thumbnails.DeleteByFileID(ctx, file.ID); err != nil {
		return fmt.Errorf("failed to delete thumbnail metadata: %w", err)
	}

	// Delete thumbnail from storage if exists
	if file.Thumbnail != nil {
		s.storage.Delete(file.Thumbnail.Hash) // Ignore errors
	}

	// Delete file from storage (ignore not found errors)
	if err := s.storage.Delete(file.Hash); err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("failed to delete file data: %w", err)
	}

	// Delete file from database
	if err := s.repo.Files.Delete(ctx, file.ID); err != nil {
		return fmt.Errorf("failed to delete file metadata: %w", err)
	}

	return nil
}

// verifyFileHash verifies that the stored file contents match the claimed SHA-256 hash
func (s *FileService) verifyFileHash(hash string) error {
	reader, err := s.storage.Get(hash)
	if err != nil {
		return fmt.Errorf("failed to open file for verification: %w", err)
	}
	defer reader.Close()

	h := sha256.New()
	if _, err := io.Copy(h, reader); err != nil {
		return fmt.Errorf("failed to read file for verification: %w", err)
	}

	calculatedHash := fmt.Sprintf("%x", h.Sum(nil))
	if calculatedHash != hash {
		return ErrHashMismatch
	}

	return nil
}

// runProcessors runs all registered processors on the file
func (s *FileService) runProcessors(ctx context.Context, file *domain.File) error {
	for _, p := range s.processors {
		if err := p.Process(ctx, file, s.storage, s.repo); err != nil {
			return fmt.Errorf("processor %s failed: %w", p.Name(), err)
		}
	}
	return nil
}

// Helper functions

func validateCreateFileRequest(req domain.CreateFileRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) > maxNameLen {
		return ErrNameTooLong
	}
	if req.Hash == "" {
		return errors.New("hash is required")
	}
	if len(req.Hash) != sha256HexLen || !sha256HexRegex.MatchString(req.Hash) {
		return ErrInvalidHash
	}
	if req.Size <= 0 {
		return errors.New("size must be greater than 0")
	}
	if req.ContentType == "" {
		return errors.New("content type is required")
	}
	if len(req.ContentType) > maxContentTypeLen {
		return ErrContentTypeTooLong
	}
	return nil
}

// maxBytesReader reads at most max bytes from r; further reads return ErrFileTooLarge
type maxBytesReader struct {
	r   io.Reader
	n   int64
	max int64
}

func newMaxBytesReader(r io.Reader, max int64) *maxBytesReader {
	return &maxBytesReader{r: r, max: max}
}

func (m *maxBytesReader) Read(p []byte) (n int, err error) {
	if m.n >= m.max {
		return 0, ErrFileTooLarge
	}
	limit := int(m.max - m.n)
	if len(p) > limit {
		p = p[:limit]
	}
	n, err = m.r.Read(p)
	m.n += int64(n)
	return n, err
}

func ptrInt32(x int32) *int32 { return &x }

// dbFileToDoamin converts a repository file to a domain file
func dbFileToDoamin(f *repository.File) *domain.File {
	return &domain.File{
		ID:            f.ID,
		Size:          f.Size,
		Name:          f.Name,
		Alias:         f.Alias,
		Hash:          f.Hash,
		Slug:          f.Slug,
		ContentType:   f.ContentType,
		CreatedAt:     timeFromPgType(f.CreatedAt),
		UpdatedAt:     timeFromPgType(f.UpdatedAt),
		UserID:        f.UserID,
		Private:       f.Private,
		Comment:       f.Comment,
		BytesReceived: f.BytesReceived,
	}
}

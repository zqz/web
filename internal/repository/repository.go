package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides access to all data repositories
type Repository struct {
	Files      FileRepository
	Users      UserRepository
	Thumbnails ThumbnailRepository
	Settings   SettingsRepository
}

// NewRepository creates a new Repository with all sub-repositories
func NewRepository(pool *pgxpool.Pool) *Repository {
	queries := New(pool)

	return &Repository{
		Files:      NewFileRepository(queries),
		Users:      NewUserRepository(queries),
		Thumbnails: NewThumbnailRepository(queries),
		Settings:   NewSettingsRepository(queries),
	}
}

// FileRepository defines the interface for file data access
type FileRepository interface {
	Create(ctx context.Context, params CreateFileParams) (*File, error)
	GetByID(ctx context.Context, id int32) (*File, error)
	GetBySlug(ctx context.Context, slug string) (*File, error)
	GetByHash(ctx context.Context, hash string) (*File, error)
	GetWithThumbnail(ctx context.Context, id int32) (*FileWithThumbnail, error)
	GetWithThumbnailBySlug(ctx context.Context, slug string) (*FileWithThumbnail, error)
	GetWithThumbnailByHash(ctx context.Context, hash string) (*FileWithThumbnail, error)
	List(ctx context.Context, limit, offset int32) ([]*File, error)
	ListByUserID(ctx context.Context, userID int32, limit, offset int32) ([]*File, error)
	ListPublic(ctx context.Context, limit, offset int32) ([]*File, error)
	ListVisibleToUser(ctx context.Context, userID int32, limit, offset int32) ([]*File, error)
	SearchFiles(ctx context.Context, search string, limit, offset int32) ([]*File, error)
	SearchPublicFiles(ctx context.Context, search string, limit, offset int32) ([]*File, error)
	SearchFilesVisibleToUser(ctx context.Context, userID int32, search string, limit, offset int32) ([]*File, error)
	ListWithThumbnails(ctx context.Context, limit, offset int32) ([]*FileWithThumbnail, error)
	Update(ctx context.Context, params UpdateFileParams) (*File, error)
	Delete(ctx context.Context, id int32) error
	DeleteByUserID(ctx context.Context, userID int32) error
	Count(ctx context.Context) (int64, error)
	CountByUserID(ctx context.Context, userID int32) (int64, error)
	TotalSize(ctx context.Context) (int64, error)
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, params CreateUserParams) (*User, error)
	GetByID(ctx context.Context, id int32) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByProviderID(ctx context.Context, providerID string) (*User, error)
	List(ctx context.Context, limit, offset int32) ([]*User, error)
	Update(ctx context.Context, params UpdateUserParams) (*User, error)
	UpdateProfile(ctx context.Context, params UpdateUserProfileParams) (*User, error)
	SetBanned(ctx context.Context, userID int32, banned bool) (*User, error)
	SetMaxFileSize(ctx context.Context, userID int32, maxBytes *int64) (*User, error)
	Delete(ctx context.Context, id int32) error
	Count(ctx context.Context) (int64, error)
	CountBanned(ctx context.Context) (int64, error)
}

// SettingsRepository defines the interface for site settings
type SettingsRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

// ThumbnailRepository defines the interface for thumbnail data access
type ThumbnailRepository interface {
	Create(ctx context.Context, params CreateThumbnailParams) (*Thumbnail, error)
	GetByFileID(ctx context.Context, fileID int32) (*Thumbnail, error)
	ListByFileID(ctx context.Context, fileID int32) ([]*Thumbnail, error)
	Update(ctx context.Context, params UpdateThumbnailParams) (*Thumbnail, error)
	Delete(ctx context.Context, id int32) error
	DeleteByFileID(ctx context.Context, fileID int32) error
}

// FileWithThumbnail represents a file with its thumbnail information
type FileWithThumbnail struct {
	File
	ThumbnailHash   *string
	ThumbnailWidth  *int32
	ThumbnailHeight *int32
}

// Transactor defines the interface for database transaction operations
type Transactor interface {
	// WithTransaction executes a function within a database transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context, repo *Repository) error) error
}

// transactor implements the Transactor interface
type transactor struct {
	pool *pgxpool.Pool
}

// NewTransactor creates a new Transactor
func NewTransactor(pool *pgxpool.Pool) Transactor {
	return &transactor{pool: pool}
}

// WithTransaction executes the given function within a transaction
func (t *transactor) WithTransaction(ctx context.Context, fn func(ctx context.Context, repo *Repository) error) error {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	// Create a repository using the transaction
	queries := New(tx)
	repo := &Repository{
		Files:      NewFileRepository(queries),
		Users:      NewUserRepository(queries),
		Thumbnails: NewThumbnailRepository(queries),
		Settings:   NewSettingsRepository(queries),
	}

	return fn(ctx, repo)
}

// Helper functions to convert between domain models and database models

// ErrNotFound is returned when a record is not found
var ErrNotFound = sql.ErrNoRows

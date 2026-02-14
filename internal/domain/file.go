package domain

import (
	"time"
)

// File represents a file in the system (domain model)
type File struct {
	ID          int32
	Size        int32
	Name        string
	Alias       string
	Hash        string
	Slug        string
	ContentType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	UserID      *int32
	Private     bool
	Comment     string

	// Additional fields not in DB
	BytesReceived int32
	Thumbnail     *Thumbnail
}

// Finished returns true if the file upload is complete
func (f *File) Finished() bool {
	return f.Size == f.BytesReceived
}

// IsOwnedBy checks if the file is owned by the given user ID
func (f *File) IsOwnedBy(userID int32) bool {
	return f.UserID != nil && *f.UserID == userID
}

// IsPublic returns true if the file is not private
func (f *File) IsPublic() bool {
	return !f.Private
}

// CanBeAccessedBy checks if a file can be accessed by a user
// Public files can be accessed by anyone
// Private files can only be accessed by the owner
func (f *File) CanBeAccessedBy(userID *int32) bool {
	if f.IsPublic() {
		return true
	}

	if userID == nil {
		return false
	}

	return f.IsOwnedBy(*userID)
}

// Thumbnail represents a thumbnail for a file
type Thumbnail struct {
	ID        int32
	FileID    int32
	Hash      string
	Width     int32
	Height    int32
	CreatedAt time.Time
}

// CreateFileRequest represents a request to create a file
type CreateFileRequest struct {
	Name        string
	Hash        string
	Size        int32
	ContentType string
	UserID      *int32
	Private     bool
	Comment     string
}

// UpdateFileRequest represents a request to update a file
type UpdateFileRequest struct {
	Name    *string
	Private *bool
	Comment *string
}

// FileUploadOptions contains options for file upload
type FileUploadOptions struct {
	GenerateThumbnail bool
	MaxThumbnailSize  int
}

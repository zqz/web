package storage

import (
	"errors"
	"io"
)

var (
	// ErrNotFound is returned when a file is not found in storage
	ErrNotFound = errors.New("file not found in storage")

	// ErrAlreadyExists is returned when trying to create a file that already exists
	ErrAlreadyExists = errors.New("file already exists in storage")
)

// Storage defines the interface for file storage operations
type Storage interface {
	// Put stores data with the given key
	// Returns an error if the file already exists
	Put(key string, data io.Reader) error

	// Append appends data to an existing file or creates it if it doesn't exist
	Append(key string, data io.Reader) (int64, error)

	// Get retrieves data for the given key
	// Returns ErrNotFound if the key doesn't exist
	Get(key string) (io.ReadCloser, error)

	// Delete removes the file with the given key
	// Returns ErrNotFound if the key doesn't exist
	Delete(key string) error

	// Exists checks if a file exists
	Exists(key string) (bool, error)

	// Size returns the size of the file in bytes
	// Returns ErrNotFound if the key doesn't exist
	Size(key string) (int64, error)

	// Path returns the full path to the file (for disk storage)
	// May return empty string for non-disk storage implementations
	Path(key string) string
}

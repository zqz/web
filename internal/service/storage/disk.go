package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DiskStorage implements Storage interface using the local filesystem
type DiskStorage struct {
	basePath string
}

// NewDiskStorage creates a new disk storage instance
func NewDiskStorage(basePath string) (*DiskStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &DiskStorage{
		basePath: basePath,
	}, nil
}

// Put stores data with the given key
func (d *DiskStorage) Put(key string, data io.Reader) error {
	path := d.fullPath(key)

	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return ErrAlreadyExists
	}

	// Create file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data
	if _, err := io.Copy(file, data); err != nil {
		os.Remove(path) // Clean up on error
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Append appends data to an existing file or creates it if it doesn't exist
func (d *DiskStorage) Append(key string, data io.Reader) (int64, error) {
	path := d.fullPath(key)

	// Open file for append, create if doesn't exist
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open file for append: %w", err)
	}
	defer file.Close()

	// Copy data and return bytes written
	n, err := io.Copy(file, data)
	if err != nil {
		return n, fmt.Errorf("failed to append to file: %w", err)
	}

	return n, nil
}

// Get retrieves data for the given key
func (d *DiskStorage) Get(key string) (io.ReadCloser, error) {
	path := d.fullPath(key)

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes the file with the given key
func (d *DiskStorage) Delete(key string) error {
	path := d.fullPath(key)

	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Exists checks if a file exists
func (d *DiskStorage) Exists(key string) (bool, error) {
	path := d.fullPath(key)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	return true, nil
}

// Size returns the size of the file in bytes
func (d *DiskStorage) Size(key string) (int64, error) {
	path := d.fullPath(key)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}

// Path returns the full path to the file
func (d *DiskStorage) Path(key string) string {
	return d.fullPath(key)
}

// fullPath returns the full filesystem path for a given key
func (d *DiskStorage) fullPath(key string) string {
	return filepath.Join(d.basePath, key)
}

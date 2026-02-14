package processor

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service/storage"
)

// ThumbnailProcessor generates thumbnails for image files
type ThumbnailProcessor struct {
	maxSize int
}

// NewThumbnailProcessor creates a new thumbnail processor
func NewThumbnailProcessor(maxSize int) *ThumbnailProcessor {
	return &ThumbnailProcessor{
		maxSize: maxSize,
	}
}

// Name returns the processor name
func (p *ThumbnailProcessor) Name() string {
	return "thumbnail"
}

// Process generates a thumbnail for image files
func (p *ThumbnailProcessor) Process(ctx context.Context, file *domain.File, stor storage.Storage, repo *repository.Repository) error {
	// Only process image files
	if !isImageContentType(file.ContentType) {
		return nil // Skip non-images
	}

	// Get file data
	reader, err := stor.Get(file.Hash)
	if err != nil {
		return fmt.Errorf("failed to get file data: %w", err)
	}
	defer reader.Close()

	// Decode image
	img, format, err := image.Decode(reader)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	// Generate thumbnail
	thumbnail := imaging.Fit(img, p.maxSize, p.maxSize, imaging.Lanczos)

	// Encode thumbnail to buffer (JPEG for photos, PNG for transparency)
	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: 85}
	switch format {
	case "png":
		err = png.Encode(&buf, thumbnail)
	default:
		err = jpeg.Encode(&buf, thumbnail, opts)
	}
	if err != nil {
		return fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	// Calculate thumbnail hash (SHA-256)
	sum := sha256.Sum256(buf.Bytes())
	thumbnailHash := fmt.Sprintf("%x", sum[:])

	// Store thumbnail
	if err := stor.Put(thumbnailHash, bytes.NewReader(buf.Bytes())); err != nil {
		// Ignore already exists errors
		if err != storage.ErrAlreadyExists {
			return fmt.Errorf("failed to store thumbnail: %w", err)
		}
	}

	// Delete existing thumbnails for this file
	if err := repo.Thumbnails.DeleteByFileID(ctx, file.ID); err != nil {
		return fmt.Errorf("failed to delete old thumbnails: %w", err)
	}

	// Save thumbnail metadata
	_, err = repo.Thumbnails.Create(ctx, repository.CreateThumbnailParams{
		FileID: file.ID,
		Hash:   thumbnailHash,
		Width:  int32(thumbnail.Bounds().Dx()),
		Height: int32(thumbnail.Bounds().Dy()),
	})
	if err != nil {
		return fmt.Errorf("failed to save thumbnail metadata: %w", err)
	}

	return nil
}

// isImageContentType checks if a content type is an image
func isImageContentType(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.HasPrefix(contentType, "image/")
}

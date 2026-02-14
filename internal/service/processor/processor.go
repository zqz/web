package processor

import (
	"context"

	"github.com/zqz/web/backend/internal/domain"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service/storage"
)

// Processor defines the interface for file processing operations
// Processors run after a file has been fully uploaded
type Processor interface {
	// Process processes a file and may create additional artifacts (thumbnails, metadata, etc.)
	Process(ctx context.Context, file *domain.File, storage storage.Storage, repo *repository.Repository) error

	// Name returns the processor name for logging
	Name() string
}

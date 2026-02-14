package repository

import (
	"context"
	"database/sql"
)

type thumbnailRepository struct {
	queries *Queries
}

// NewThumbnailRepository creates a new thumbnail repository
func NewThumbnailRepository(queries *Queries) ThumbnailRepository {
	return &thumbnailRepository{queries: queries}
}

func (r *thumbnailRepository) Create(ctx context.Context, params CreateThumbnailParams) (*Thumbnail, error) {
	thumbnail, err := r.queries.CreateThumbnail(ctx, params)
	if err != nil {
		return nil, err
	}
	return &thumbnail, nil
}

func (r *thumbnailRepository) GetByFileID(ctx context.Context, fileID int32) (*Thumbnail, error) {
	thumbnail, err := r.queries.GetThumbnailByFileID(ctx, fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &thumbnail, nil
}

func (r *thumbnailRepository) ListByFileID(ctx context.Context, fileID int32) ([]*Thumbnail, error) {
	thumbnails, err := r.queries.GetThumbnailsByFileID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	result := make([]*Thumbnail, len(thumbnails))
	for i := range thumbnails {
		result[i] = &thumbnails[i]
	}
	return result, nil
}

func (r *thumbnailRepository) Update(ctx context.Context, params UpdateThumbnailParams) (*Thumbnail, error) {
	thumbnail, err := r.queries.UpdateThumbnail(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &thumbnail, nil
}

func (r *thumbnailRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteThumbnail(ctx, id)
}

func (r *thumbnailRepository) DeleteByFileID(ctx context.Context, fileID int32) error {
	return r.queries.DeleteThumbnailsByFileID(ctx, fileID)
}

package repository

import (
	"context"
	"database/sql"
)

type fileRepository struct {
	queries *Queries
}

// NewFileRepository creates a new file repository
func NewFileRepository(queries *Queries) FileRepository {
	return &fileRepository{queries: queries}
}

func (r *fileRepository) Create(ctx context.Context, params CreateFileParams) (*File, error) {
	file, err := r.queries.CreateFile(ctx, params)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetByID(ctx context.Context, id int32) (*File, error) {
	file, err := r.queries.GetFileByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetBySlug(ctx context.Context, slug string) (*File, error) {
	file, err := r.queries.GetFileBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetByHash(ctx context.Context, hash string) (*File, error) {
	file, err := r.queries.GetFileByHash(ctx, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) GetWithThumbnail(ctx context.Context, id int32) (*FileWithThumbnail, error) {
	row, err := r.queries.GetFileWithThumbnail(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return rowToFileWithThumbnail(row), nil
}

func (r *fileRepository) GetWithThumbnailBySlug(ctx context.Context, slug string) (*FileWithThumbnail, error) {
	row, err := r.queries.GetFileWithThumbnailBySlug(ctx, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &FileWithThumbnail{
		File: File{
			ID:          row.ID,
			Size:        row.Size,
			Name:        row.Name,
			Alias:       row.Alias,
			Hash:        row.Hash,
			Slug:        row.Slug,
			ContentType: row.ContentType,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			UserID:      row.UserID,
			Private:     row.Private,
			Comment:     row.Comment,
		},
		ThumbnailHash:   row.ThumbnailHash,
		ThumbnailWidth:  row.ThumbnailWidth,
		ThumbnailHeight: row.ThumbnailHeight,
	}, nil
}

func (r *fileRepository) GetWithThumbnailByHash(ctx context.Context, hash string) (*FileWithThumbnail, error) {
	row, err := r.queries.GetFileWithThumbnailByHash(ctx, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &FileWithThumbnail{
		File: File{
			ID:          row.ID,
			Size:        row.Size,
			Name:        row.Name,
			Alias:       row.Alias,
			Hash:        row.Hash,
			Slug:        row.Slug,
			ContentType: row.ContentType,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			UserID:      row.UserID,
			Private:     row.Private,
			Comment:     row.Comment,
		},
		ThumbnailHash:   row.ThumbnailHash,
		ThumbnailWidth:  row.ThumbnailWidth,
		ThumbnailHeight: row.ThumbnailHeight,
	}, nil
}

func (r *fileRepository) List(ctx context.Context, limit, offset int32) ([]*File, error) {
	files, err := r.queries.ListFiles(ctx, ListFilesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) ListByUserID(ctx context.Context, userID int32, limit, offset int32) ([]*File, error) {
	files, err := r.queries.ListFilesByUserID(ctx, ListFilesByUserIDParams{
		UserID: &userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) ListPublic(ctx context.Context, limit, offset int32) ([]*File, error) {
	files, err := r.queries.ListPublicFiles(ctx, ListPublicFilesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) ListVisibleToUser(ctx context.Context, userID int32, limit, offset int32) ([]*File, error) {
	files, err := r.queries.ListFilesVisibleToUser(ctx, ListFilesVisibleToUserParams{
		UserID: &userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) SearchFiles(ctx context.Context, search string, limit, offset int32) ([]*File, error) {
	files, err := r.queries.SearchFiles(ctx, SearchFilesParams{
		Name:   search,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) SearchPublicFiles(ctx context.Context, search string, limit, offset int32) ([]*File, error) {
	files, err := r.queries.SearchPublicFiles(ctx, SearchPublicFilesParams{
		Name:   search,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) SearchFilesVisibleToUser(ctx context.Context, userID int32, search string, limit, offset int32) ([]*File, error) {
	files, err := r.queries.SearchFilesVisibleToUser(ctx, SearchFilesVisibleToUserParams{
		UserID: &userID,
		Name:   search,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]*File, len(files))
	for i := range files {
		result[i] = &files[i]
	}
	return result, nil
}

func (r *fileRepository) ListWithThumbnails(ctx context.Context, limit, offset int32) ([]*FileWithThumbnail, error) {
	rows, err := r.queries.ListFilesWithThumbnails(ctx, ListFilesWithThumbnailsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*FileWithThumbnail, len(rows))
	for i, row := range rows {
		result[i] = &FileWithThumbnail{
			File: File{
				ID:          row.ID,
				Size:        row.Size,
				Name:        row.Name,
				Alias:       row.Alias,
				Hash:        row.Hash,
				Slug:        row.Slug,
				ContentType: row.ContentType,
				CreatedAt:   row.CreatedAt,
				UpdatedAt:   row.UpdatedAt,
				UserID:      row.UserID,
				Private:     row.Private,
				Comment:     row.Comment,
			},
			ThumbnailHash:   row.ThumbnailHash,
			ThumbnailWidth:  row.ThumbnailWidth,
			ThumbnailHeight: row.ThumbnailHeight,
		}
	}
	return result, nil
}

func (r *fileRepository) Update(ctx context.Context, params UpdateFileParams) (*File, error) {
	file, err := r.queries.UpdateFile(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteFile(ctx, id)
}

func (r *fileRepository) DeleteByUserID(ctx context.Context, userID int32) error {
	return r.queries.DeleteFilesByUserID(ctx, &userID)
}

func (r *fileRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountFiles(ctx)
}

func (r *fileRepository) CountByUserID(ctx context.Context, userID int32) (int64, error) {
	return r.queries.CountFilesByUserID(ctx, &userID)
}

func (r *fileRepository) TotalSize(ctx context.Context) (int64, error) {
	return r.queries.TotalFileSize(ctx)
}

// Helper function to convert GetFileWithThumbnailRow to FileWithThumbnail
func rowToFileWithThumbnail(row GetFileWithThumbnailRow) *FileWithThumbnail {
	return &FileWithThumbnail{
		File: File{
			ID:          row.ID,
			Size:        row.Size,
			Name:        row.Name,
			Alias:       row.Alias,
			Hash:        row.Hash,
			Slug:        row.Slug,
			ContentType: row.ContentType,
			CreatedAt:   row.CreatedAt,
			UpdatedAt:   row.UpdatedAt,
			UserID:      row.UserID,
			Private:     row.Private,
			Comment:     row.Comment,
		},
		ThumbnailHash:   row.ThumbnailHash,
		ThumbnailWidth:  row.ThumbnailWidth,
		ThumbnailHeight: row.ThumbnailHeight,
	}
}

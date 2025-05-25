package filedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zqz/web/backend/models"
)

type DBMetaStorage struct {
	entries      map[string]*Meta
	entriesMutex sync.Mutex
	db           *sql.DB
	ctx          context.Context
}

func file2meta(f *models.File) Meta {
	m := Meta{
		ID:            f.ID,
		UserID:        f.UserID.Int,
		Name:          f.Name,
		Size:          f.Size,
		BytesReceived: f.Size,
		Slug:          f.Slug,
		Hash:          f.Hash,
		ContentType:   f.ContentType,
		Date:          f.CreatedAt.Time,
		Comment:       f.Comment,
		Private:       f.Private,
	}

	return m
}

func meta2file(m *Meta) *models.File {
	return &models.File{
		ID:          m.ID,
		Name:        m.Name,
		Alias:       m.Name,
		Size:        m.Size,
		Slug:        m.Slug,
		Hash:        m.Hash,
		ContentType: m.ContentType,
		Private:     m.Private,
		Comment:     m.Comment,
		UserID:      null.IntFrom(m.UserID),
	}
}

const paginationSQL = `
	SELECT
	f.id, f.hash, f.name, f.slug, f.content_type, f.created_at, f.size, f.private, f.comment, t.hash as thash
	FROM files AS f
	LEFT JOIN thumbnails t ON t.file_id = f.id
	ORDER BY f.created_at DESC
	OFFSET $1
	LIMIT $2
`

func (m *DBMetaStorage) ListFilesByUserId(userID, page int) ([]*Meta, error) {
	var err error

	//	var perPage int = 50
	// offset := perPage * page

	files, err := models.Files(
		qm.Load(models.FileRels.Thumbnails),
		qm.Where("user_id=?", userID),
	).All(m.ctx, m.db)

	if err != nil {
		return nil, err
	}

	metas := make([]*Meta, 0)

	for _, f := range files {
		thumbs := f.R.Thumbnails
		thumbHash := ""
		if len(thumbs) > 0 {
			thumbHash = thumbs[0].Hash
		}

		metas = append(metas, &Meta{
			ID:        f.ID,
			Name:      f.Name,
			Slug:      f.Slug,
			Hash:      f.Hash,
			Size:      f.Size,
			Thumbnail: thumbHash,
		})
	}

	return metas, nil
}

func (m *DBMetaStorage) ListPage(page int) ([]*Meta, error) {
	entries := make([]*Meta, 0)
	var rows *sql.Rows
	var err error

	var perPage int = 50
	offset := perPage * page

	if rows, err = m.db.Query(paginationSQL, offset, perPage); err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		var e struct {
			ID            int
			Hash          null.String
			Name          null.String
			Slug          null.String
			ContentType   null.String
			Date          null.Time
			Private       bool
			Comment       string
			Size          int
			ThumbnailHash null.String
		}

		err = rows.Scan(
			&e.ID, &e.Hash, &e.Name, &e.Slug, &e.ContentType, &e.Date, &e.Size, &e.Private, &e.Comment, &e.ThumbnailHash,
		)

		if err != nil {
			fmt.Println("err", err.Error())
		}

		x := Meta{}

		x.ID = e.ID

		if e.Hash.Valid {
			x.Hash = e.Hash.String
		}

		if e.Name.Valid {
			x.Name = e.Name.String
		}

		if e.Slug.Valid {
			x.Slug = e.Slug.String
		}

		if e.ContentType.Valid {
			x.ContentType = e.ContentType.String
		}

		if e.Date.Valid {
			x.Date = e.Date.Time
		}

		if e.ThumbnailHash.Valid {
			x.Thumbnail = e.ThumbnailHash.String
		}

		x.Size = e.Size
		x.Private = e.Private
		x.Comment = e.Comment

		entries = append(entries, &x)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return entries, err
}

func NewDBMetaStorage(db *sql.DB) *DBMetaStorage {
	return &DBMetaStorage{
		entries: make(map[string]*Meta, 0),
		db:      db,
		ctx:     context.TODO(),
	}
}

func (s *DBMetaStorage) fetchMetaFromDBWithHash(h string) (Meta, error) {
	f, err := models.Files(qm.Where("hash=?", h)).One(s.ctx, s.db)

	if err != nil {
		return Meta{}, err
	}

	m := file2meta(f)

	t, err := f.Thumbnails().One(s.ctx, s.db)
	if err != nil {
		return m, nil
	}

	m.Thumbnail = t.Hash
	return m, nil
}

func (s *DBMetaStorage) fetchMetaFromDBWithSlug(slug string) (Meta, error) {
	f, err := models.Files(qm.Where("slug=?", slug)).One(s.ctx, s.db)
	m := file2meta(f)

	t, err := f.Thumbnails().One(s.ctx, s.db)
	if err != nil {
		return m, nil
	}

	m.Thumbnail = t.Hash
	return m, nil
}

func (s *DBMetaStorage) RemoveThumbnails(m *Meta) error {
	_, err := models.Thumbnails(qm.Where("file_id=?", m.ID)).DeleteAll(s.ctx, s.db)
	return err
}

func (s *DBMetaStorage) StoreThumbnail(h string, size int, m *Meta) error {
	err := s.RemoveThumbnails(m)
	if err != nil {
		return err
	}

	t := models.Thumbnail{
		Hash:   h,
		Width:  size,
		Height: size,
		FileID: m.ID,
	}

	err = t.Insert(s.ctx, s.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func (s *DBMetaStorage) FetchMeta(h string) (*Meta, error) {
	m, ok := s.entries[h]

	if ok {
		return m, nil
	}

	m2, err := s.fetchMetaFromDBWithHash(h)
	if err != nil {
		return nil, errors.New("file not found")
	}

	s.entriesMutex.Lock()
	s.entries[m2.Hash] = &m2
	s.entriesMutex.Unlock()

	return &m2, nil
}

func (s *DBMetaStorage) FetchMetaWithSlug(slug string) (*Meta, error) {
	m, err := s.fetchMetaFromDBWithSlug(slug)
	if err != nil {
		return nil, errors.New("file not found")
	}

	s.entriesMutex.Lock()
	s.entries[m.Hash] = &m
	s.entriesMutex.Unlock()

	return &m, nil
}

func (s *DBMetaStorage) UpdateMeta(m *Meta) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	f := meta2file(m)
	_, err := f.Update(s.ctx, s.db, boil.Infer())

	return err
}

func (s *DBMetaStorage) StoreMeta(m *Meta) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	if m.Finished() {
		f := meta2file(m)
		if err := f.Insert(s.ctx, s.db, boil.Infer()); err != nil {
			fmt.Println("error", err.Error())
			return err
		}

		m.ID = f.ID
	}

	return nil
}

func (s *DBMetaStorage) DeleteMetaById(id int) error {
	f, err := models.Files(qm.Where("id=?", id)).One(s.ctx, s.db)
	if err != nil {
		return err
	}

	_, err = f.Delete(s.ctx, s.db)

	return err
}

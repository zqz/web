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
	"github.com/zqz/upl/models"
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
		Name:          f.Name,
		Size:          f.Size,
		BytesReceived: f.Size,
		Slug:          f.Slug,
		Hash:          f.Hash,
		ContentType:   f.ContentType,
		Date:          f.CreatedAt.Time,
	}

	return m
}

func meta2file(m *Meta) *models.File {
	return &models.File{
		Name:        m.Name,
		Alias:       m.Name,
		Size:        m.Size,
		Slug:        m.Slug,
		Hash:        m.Hash,
		ContentType: m.ContentType,
	}
}

const paginationSQL = `
	SELECT
	f.id, f.hash, f.name, f.slug, f.content_type, f.created_at, f.size
	FROM files AS f
	ORDER BY f.created_at DESC
	OFFSET $1
	LIMIT $2
`

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
			ID          int
			Hash        null.String
			Name        null.String
			Slug        null.String
			ContentType null.String
			Date        null.Time
			Size        int
		}

		err = rows.Scan(
			&e.ID, &e.Hash, &e.Name, &e.Slug, &e.ContentType, &e.Date, &e.Size,
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
		x.Size = e.Size

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

	return file2meta(f), nil
}

func (s *DBMetaStorage) fetchMetaFromDBWithSlug(slug string) (Meta, error) {
	f, err := models.Files(qm.Where("slug=?", slug)).One(s.ctx, s.db)

	if err != nil {
		return Meta{}, err
	}

	return file2meta(f), nil
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

func (s *DBMetaStorage) StoreMeta(m *Meta) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	if m.finished() {
		f := meta2file(m)
		if err := f.Insert(s.ctx, s.db, boil.Infer()); err != nil {
			fmt.Println("error", err.Error())
			return err
		}

		m.ID = f.ID
	}

	return nil
}

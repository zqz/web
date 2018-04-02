package filedb

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/zqz/upl/models"
	"gopkg.in/volatiletech/null.v6"
)

type DBMetaStorage struct {
	entries      map[string]*Meta
	entriesMutex sync.Mutex
	db           *sql.DB
}

func file2meta(f *models.File) Meta {
	m := Meta{
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
	f.id, f.hash, f.name, f.slug, f.created_at, f.size
	FROM files AS f
	ORDER BY f.created_at DESC
	OFFSET $1
	LIMIT $2
`

func (m DBMetaStorage) ListPage(page int) ([]*Meta, error) {
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
			ID   int
			Name null.String
			Hash null.String
			Slug null.String
			Date null.Time
			Size int
		}

		err = rows.Scan(
			&e.ID, &e.Hash, &e.Name, &e.Slug, &e.Date, &e.Size,
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

		x.Size = e.Size

		entries = append(entries, &x)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return entries, err
}

func NewDBMetaStorage(db *sql.DB) DBMetaStorage {
	return DBMetaStorage{
		entries: make(map[string]*Meta, 0),
		db:      db,
	}
}

func (m DBMetaStorage) fetchMetaFromDBWithHash(hash string) (Meta, error) {
	f, err := models.Files(m.db, qm.Where("hash=?", hash)).One()

	if err != nil {
		return Meta{}, err
	}

	return file2meta(f), nil
}

func (m DBMetaStorage) fetchMetaFromDBWithSlug(slug string) (Meta, error) {
	f, err := models.Files(m.db, qm.Where("slug=?", slug)).One()
	if err != nil {
		return Meta{}, err
	}

	return file2meta(f), nil
}

func (m DBMetaStorage) FetchMeta(hash string) (*Meta, error) {
	meta, ok := m.entries[hash]

	if ok {
		return meta, nil
	}

	meta2, err := m.fetchMetaFromDBWithHash(hash)
	if err != nil {
		return nil, errors.New("file not found")
	}

	m.entriesMutex.Lock()
	m.entries[meta2.Hash] = &meta2
	m.entriesMutex.Unlock()

	return &meta2, nil
}

func (m DBMetaStorage) FetchMetaWithSlug(slug string) (*Meta, error) {
	meta, err := m.fetchMetaFromDBWithSlug(slug)
	if err != nil {
		return nil, errors.New("file not found")
	}

	m.entriesMutex.Lock()
	m.entries[meta.Hash] = &meta
	m.entriesMutex.Unlock()

	return &meta, nil
}

func (s DBMetaStorage) StoreMeta(m *Meta) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	if m.finished() {
		f := meta2file(m)
		if err := f.Insert(s.db); err != nil {
			fmt.Println("error", err.Error())
			return err
		}

		m.ID = f.ID
	}

	return nil
}

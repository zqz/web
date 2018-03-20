package filedb

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/zqz/upl/models"
)

type DBMetaStorage struct {
	entries      map[string]*Meta
	entriesMutex sync.Mutex
	db           *sql.DB
}

func file2meta(f *models.File) *Meta {
	return &Meta{
		Name:          f.Name,
		Size:          f.Size,
		BytesReceived: f.Size,
		Slug:          f.Slug,
		Hash:          f.Hash,
		ContentType:   f.ContentType,
		Date:          f.CreatedAt.Time,
	}
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

func (m DBMetaStorage) ListPage(page int) ([]*Meta, error) {
	metas := make([]*Meta, 0)

	files, err := models.Files(m.db, qm.Limit(10)).All()

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		metas = append(metas, file2meta(f))
	}

	return metas, nil
}

func NewDBMetaStorage(db *sql.DB) DBMetaStorage {
	return DBMetaStorage{
		entries: make(map[string]*Meta, 0),
		db:      db,
	}
}

func (m DBMetaStorage) fetchMetaFromDBWithHash(hash string) (*Meta, error) {
	f, err := models.Files(m.db, qm.Where("hash=?", hash)).One()
	if err != nil {
		return nil, err
	}

	return file2meta(f), nil
}

func (m DBMetaStorage) fetchMetaFromDBWithSlug(slug string) (*Meta, error) {
	f, err := models.Files(m.db, qm.Where("slug=?", slug)).One()
	if err != nil {
		return nil, err
	}

	return file2meta(f), nil
}

func (m DBMetaStorage) FetchMeta(hash string) (*Meta, error) {
	meta, ok := m.entries[hash]

	if ok {
		return meta, nil
	}

	meta, err := m.fetchMetaFromDBWithHash(hash)
	if err != nil {
		return nil, errors.New("file not found")
	}

	m.entriesMutex.Lock()
	m.entries[meta.Hash] = meta
	m.entriesMutex.Unlock()

	return meta, nil
}

func (m DBMetaStorage) FetchMetaWithSlug(slug string) (*Meta, error) {
	meta, err := m.fetchMetaFromDBWithSlug(slug)
	if err != nil {
		return nil, errors.New("file not found")
	}

	m.entriesMutex.Lock()
	m.entries[meta.Hash] = meta
	m.entriesMutex.Unlock()

	return meta, nil
}

func (m DBMetaStorage) StoreMeta(meta Meta) error {
	m.entriesMutex.Lock()
	m.entries[meta.Hash] = &meta
	m.entriesMutex.Unlock()

	if meta.finished() {
		f := meta2file(&meta)
		err := f.Insert(m.db)
		if err != nil {
			fmt.Println("error", err.Error())
			return err
		}
	}

	return nil
}

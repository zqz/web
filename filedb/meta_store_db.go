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
	t.hash, f.hash, f.name, f.slug, f.created_at, f.size
	FROM files AS f
	LEFT JOIN thumbnails as t
	ON t.file_id = f.id
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
			Thumbnail null.String
			Name      null.String
			Hash      null.String
			Slug      null.String
			Date      null.Time
			Size      int
		}

		err = rows.Scan(
			&e.Thumbnail, &e.Hash, &e.Name, &e.Slug, &e.Date, &e.Size,
		)

		if err != nil {
			fmt.Println("err", err.Error())
		}

		x := Meta{}
		if e.Thumbnail.Valid {
			x.Thumbnail = e.Thumbnail.String
		}

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

func (m DBMetaStorage) StoreThumbnail(t Thumbnail) error {
	fmt.Println("looking for", t.MetaHash)
	f, err := models.Files(m.db, qm.Where("hash=?", t.MetaHash)).One()

	if err != nil {
		fmt.Println("error fetching file", err.Error())
		return err
	}

	tn := models.Thumbnail{
		FileID: null.IntFrom(f.ID),
		Hash:   t.Hash,
	}

	err = tn.Insert(m.db)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}

	return nil
}

func (m DBMetaStorage) ThumbnailExists(h string) (bool, error) {
	_, err := models.Thumbnails(m.db, qm.Where("hash=?", h)).One()

	if err != nil {
		return false, err
	}

	return true, nil
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

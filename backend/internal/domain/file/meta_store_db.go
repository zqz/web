package file

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zqz/web/backend/internal/models"
)

type DBMetaStorage struct {
	entries      map[string]*File
	entriesMutex sync.Mutex
	db           *sql.DB
	ctx          context.Context
}

// func xxfile2meta(f *models.File) File {
// 	m := Meta{
// 		ID:            f.ID,
// 		UserID:        f.UserID.Int,
// 		Name:          f.Name,
// 		Size:          f.Size,
// 		BytesReceived: f.Size,
// 		Slug:          f.Slug,
// 		Hash:          f.Hash,
// 		ContentType:   f.ContentType,
// 		Date:          f.CreatedAt.Time,
// 		Comment:       f.Comment,
// 		Private:       f.Private,
// 	}
//
// 	return m
// }
//
// func xxmeta2file(m *File) *models.File {
// 	f := models.File{
// 		ID:          m.ID,
// 		Name:        m.Name,
// 		Alias:       m.Name,
// 		Size:        m.Size,
// 		Slug:        m.Slug,
// 		Hash:        m.Hash,
// 		ContentType: m.ContentType,
// 		Private:     m.Private,
// 		Comment:     m.Comment,
// 	}
//
// 	if m.UserID > 0 {
// 		f.UserID = null.IntFrom(m.UserID)
// 	}
//
// 	return &f
// }

const paginationSQL = `
	SELECT
	f.id, f.hash, f.name, f.slug, f.content_type, f.created_at, f.size, f.private, f.comment, t.hash as thash
	FROM files AS f
	LEFT JOIN thumbnails t ON t.file_id = f.id
	ORDER BY f.created_at DESC
	OFFSET $1
	LIMIT $2
`

func (m *DBMetaStorage) ListFilesByUserId(userID, page int) ([]*File, error) {
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

	metas := make([]*File, 0)

	for _, dbf := range files {
		thumbs := dbf.R.Thumbnails
		thumbHash := ""
		if len(thumbs) > 0 {
			thumbHash = thumbs[0].Hash
		}

		f := File{
			File:      *dbf,
			Thumbnail: thumbHash,
		}

		metas = append(metas, &f)
	}

	return metas, nil
}

func (m *DBMetaStorage) ListPage(page int) ([]*File, error) {
	entries := make([]*File, 0)
	var rows *sql.Rows
	var err error

	var perPage int = 50
	offset := perPage * page

	if rows, err = m.db.Query(paginationSQL, offset, perPage); err != nil {
		return entries, err
	}
	defer rows.Close()

	for rows.Next() {
		e := File{}
		err = rows.Scan(
			&e.ID, &e.Hash, &e.Name, &e.Slug, &e.ContentType, &e.CreatedAt, &e.Size, &e.Private, &e.Comment, &e.Thumbnail,
		)

		if err != nil {
			fmt.Println("err", err.Error())
		}

		entries = append(entries, &e)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return entries, err
}

func NewDBMetaStorage(db *sql.DB) *DBMetaStorage {
	return &DBMetaStorage{
		entries: make(map[string]*File, 0),
		db:      db,
		ctx:     context.TODO(),
	}
}

func (s *DBMetaStorage) fetchMetaFromDBWithHash(h string) (File, error) {
	f, err := models.Files(qm.Where("hash=?", h)).One(s.ctx, s.db)

	if err != nil {
		return File{}, err
	}

	m := File{File: *f}

	t, err := f.Thumbnails().One(s.ctx, s.db)
	if err != nil {
		return m, nil
	}

	m.Thumbnail = t.Hash
	return m, nil
}

func (s *DBMetaStorage) fetchMetaFromDBWithSlug(slug string) (File, error) {
	f, err := models.Files(qm.Where("slug=?", slug)).One(s.ctx, s.db)
	if f == nil {
		return File{}, errors.New("failed to find meta")
	}

	m := File{File: *f}

	t, err := f.Thumbnails().One(s.ctx, s.db)
	if err != nil {
		return m, nil
	}

	m.Thumbnail = t.Hash
	return m, nil
}

func (s *DBMetaStorage) RemoveThumbnails(m *File) error {
	_, err := models.Thumbnails(qm.Where("file_id=?", m.ID)).DeleteAll(s.ctx, s.db)
	return err
}

func (s *DBMetaStorage) StoreThumbnail(h string, size int, m *File) error {
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

func (s *DBMetaStorage) FetchByHash(h string) (*File, error) {
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

func (s *DBMetaStorage) FetchBySlug(slug string) (*File, error) {
	m, err := s.fetchMetaFromDBWithSlug(slug)
	if err != nil {
		return nil, errors.New("file not found")
	}

	s.entriesMutex.Lock()
	s.entries[m.Hash] = &m
	s.entriesMutex.Unlock()

	return &m, nil
}

func (s *DBMetaStorage) Update(m *File) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	_, err := m.Update(s.ctx, s.db, boil.Infer())

	return err
}

func (s *DBMetaStorage) Create(m *File) error {
	s.entriesMutex.Lock()
	s.entries[m.Hash] = m
	s.entriesMutex.Unlock()

	if m.Finished() {
		if err := m.Insert(s.ctx, s.db, boil.Infer()); err != nil {
			fmt.Println("error", err.Error())
			return err
		}
	}

	return nil
}

func (s *DBMetaStorage) DeleteById(id int) error {
	f, err := models.Files(qm.Where("id=?", id)).One(s.ctx, s.db)
	if err != nil {
		return err
	}

	_, err = f.Delete(s.ctx, s.db)

	return err
}

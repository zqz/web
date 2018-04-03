package filedb

import (
	"database/sql"
	"errors"

	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/zqz/upl/models"
)

type DBThumbnailStorage struct {
	db *sql.DB
}

func NewDBThumbnailStorage(db *sql.DB) DBThumbnailStorage {
	return DBThumbnailStorage{
		db: db,
	}
}

func thumb2db(t Thumbnail) models.Thumbnail {
	return models.Thumbnail{
		FileID: t.MetaID,
		Hash:   t.Hash,
		Width:  t.Width,
		Height: t.Height,
	}
}

func db2thumb(t *models.Thumbnail) Thumbnail {
	return Thumbnail{
		MetaID: t.FileID,
		Hash:   t.Hash,
		Width:  t.Width,
		Height: t.Height,
	}
}

func (s DBThumbnailStorage) StoreThumbnail(t Thumbnail) error {
	if t.MetaID == 0 {
		return errors.New("thumbnail missing required meta id")
	}

	dbt := thumb2db(t)

	return dbt.Insert(s.db)
}

func (s DBThumbnailStorage) FetchThumbnails(ids []int) (map[int]Thumbnail, error) {
	ts := make(map[int]Thumbnail, 0)

	dbIds := make([]interface{}, len(ids))
	for i, v := range ids {
		dbIds[i] = v
	}
	dbts, err := models.Thumbnails(
		s.db,
		qm.WhereIn("file_id in ?", dbIds...),
	).All()

	if err != nil {
		return ts, err
	}

	for _, dbt := range dbts {
		t := db2thumb(dbt)
		ts[t.MetaID] = t
	}

	return ts, nil
}

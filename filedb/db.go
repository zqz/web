package filedb

//import (
//	"errors"
//	"io"
//	"os"

//	"github.com/jmoiron/sqlx"
//	"github.com/volatiletech/sqlboiler/queries/qm"
//	"github.com/zqz/upl/models"
//)

//func modelToFile(mf *models.File) *File {
//	return &File{
//		Name:        mf.Name,
//		Hash:        mf.Hash,
//		Slug:        mf.Token,
//		ContentType: mf.ContentType,
//		Path:        "/tmp/final/" + mf.Hash,
//		Size:        mf.Size,
//		Date:        mf.CreatedAt.Time,
//	}
//}

//type dbFile struct {
//	*File
//	bytesReceived int
//}

//func (f dbFile) tmpPath() string {
//	return "/tmp/zqz/" + f.Hash
//}

//func (f dbFile) finalPath() string {
//	return "/tmp/final" + f.Hash
//}

//type DBFileManager struct {
//	files map[string]*dbFile
//	db    *sqlx.DB
//}

//func NewDBFileManager(c *sqlx.DB) *DBFileManager {
//	return &DBFileManager{
//		files: make(map[string]*dbFile),
//		db:    c,
//	}
//}
//func (dm *DBFileManager) FindBySlug(s string) (*File, error) {
//	if f, err := models.Files(dm.db, qm.Where("slug=?", s)).One(); err == nil {
//		return modelToFile(f), nil
//	}

//	return nil, errors.New("file not found")
//}

//func (dm *DBFileManager) FindByHash(h string) (*File, error) {
//	if f, err := models.Files(dm.db, qm.Where("hash=?", h)).One(); err == nil {
//		return modelToFile(f), nil
//	}

//	return nil, errors.New("file not found")
//}

//func (dm *DBFileManager) Finished(h string) bool {
//	df, ok := dm.files[h]

//	return ok && df.Size == df.bytesReceived
//}

//func (dm *DBFileManager) moveFile(f *dbFile) error {
//	return os.Rename(f.tmpPath(), f.finalPath())
//}

//func (dm *DBFileManager) insertFile(f *dbFile) error {
//	mf := models.File{
//		Name:        f.Name,
//		Size:        f.Size,
//		Hash:        f.Hash,
//		Token:       randStr(5),
//		ContentType: f.ContentType,
//	}

//	if err := mf.Insert(dm.db); err != nil {
//		return err
//	}

//	delete(dm.files, f.Hash)

//	return nil
//}

//func (dm *DBFileManager) StoreMeta(m Meta) (int, error) {
//	var df *dbFile

//	df, ok := dm.files[m.Hash]

//	if ok {

//		//if uploaded, return error

//		return df.bytesReceived, nil
//	}

//	df = &dbFile{
//		File: &File{
//			Name:        m.Name,
//			Size:        m.Size,
//			Hash:        m.Hash,
//			ContentType: m.ContentType,
//		},
//		bytesReceived: 0,
//	}

//	dm.files[m.Hash] = df

//	return 0, nil
//}

//func (dm *DBFileManager) StoreData(h string, r io.Reader) error {
//	if dm.Finished(h) {
//		return errors.New("file already uploaded")
//	}

//	df, ok := dm.files[h]
//	if !ok {
//		return errors.New("file not found")
//	}

//	f, err := os.OpenFile("/tmp/zqz/"+h, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
//	if err != nil {
//		return err
//	}

//	defer f.Close()

//	// ignore here. likely only ever an EOF error which is expected.
//	i, _ := io.Copy(f, r)
//	df.bytesReceived += int(i)

//	if dm.Finished(h) {
//		if err := dm.insertFile(df); err != nil {
//			return err
//		}

//		if err := dm.moveFile(df); err != nil {
//			return err
//		}
//	}

//	return nil
//}

//func (dm *DBFileManager) Write(h string, w io.Writer) error {
//	path := "/home/zqz/" + h
//	f, err := os.Open(path)

//	if err != nil {
//		return err
//	}

//	defer f.Close()

//	_, err = io.Copy(w, f)

//	return err
//}

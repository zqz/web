package file

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type persister interface {
	Del(string) error
	Put(string) (io.WriteCloser, error)
	Get(string) (io.ReadCloser, error)
	Path(string) string
}

type metaStorer interface {
	DeleteMetaById(int) error
	FetchMetaWithSlug(string) (*File, error)
	FetchMeta(string) (*File, error)
	StoreMeta(*File) error
	StoreThumbnail(string, int, *File) error
	RemoveThumbnails(*File) error
	UpdateMeta(*File) error
	ListPage(int) ([]*File, error)
	ListFilesByUserId(int, int) ([]*File, error)
}

// DB implements a upload server.
type DB struct {
	p  persister
	m  metaStorer
	px []processor
}

func (db *DB) ListFilesByUserId(uID, offset int) ([]*File, error) {
	return db.m.ListFilesByUserId(uID, offset)
}

func (db *DB) AddProcessor(p processor) {
	db.px = append(db.px, p)
}

type writeCounter int64

func (w writeCounter) Write(b []byte) (int, error) {
	w += writeCounter(len(b))

	return len(b), nil
}

// NewDB returns an instance of a DB.
func NewDB(p persister, m metaStorer) DB {
	return DB{
		m:  m,
		p:  p,
		px: make([]processor, 0),
	}
}

func (db DB) Process(m *File) error {
	err := db.process(m)
	return err
}

func (db DB) List(page int) ([]*File, error) {
	metas, err := db.m.ListPage(page)
	if err != nil {
		return nil, err
	}

	metaIds := make([]int, len(metas))
	for i, m := range metas {
		metaIds[i] = m.ID
	}

	return metas, nil
}

func (db DB) Write(hash string, rc io.ReadCloser) (*File, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	m, err := db.fetch(hash)
	if err != nil {
		return nil, err
	}

	if err := db.store(m, rc); err != nil {
		return m, err
	}

	if err := db.finish(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (db DB) Read(hash string, wc io.Writer) error {
	if err := db.validate(); err != nil {
		return err
	}

	m, err := db.m.FetchMeta(hash)
	if err != nil {
		return err
	}

	if !m.Finished() {
		return errors.New("file incomplete")
	}

	reader, err := db.p.Get(hash)

	if err != nil {
		return err
	}

	_, err = io.Copy(wc, reader)

	return err
}

func (db DB) StoreMeta(meta File) error {
	if err := db.validate(); err != nil {
		return err
	}

	if err := validateMeta(&meta); err != nil {
		return err
	}

	meta.BytesReceived = 0

	m, _ := db.m.FetchMeta(meta.Hash)
	if m != nil {
		if meta.Size != m.Size {
			return errors.New("can not change file size")
		}

		meta.BytesReceived = m.BytesReceived
	}

	if meta.BytesReceived == 0 {
		// if for some reason there is /already/ some data under this hash and we
		// have not saved any bytes. Delete the file.
		db.p.Del(meta.Hash)
	}

	return db.m.StoreMeta(&meta)
}

func (db DB) FetchMeta(h string) (*File, error) {
	return db.fetch(h)
}

func (db DB) UpdateMeta(m *File) error {
	return db.m.UpdateMeta(m)
}

func (db DB) fetch(hash string) (*File, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	if hash == "" {
		return nil, errors.New("no hash specified")
	}

	return db.m.FetchMeta(hash)
}

func (db DB) FetchMetaWithSlug(slug string) (*File, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	if slug == "" {
		return nil, errors.New("no slug specified")
	}

	return db.m.FetchMetaWithSlug(slug)
}

func (db DB) DeleteMetaWithSlug(slug string) error {
	m, err := db.m.FetchMetaWithSlug(slug)
	if err != nil {
		fmt.Println("failed to fetch meta")
		return err
	}

	err = db.p.Del(m.Thumbnail)
	if err != nil {
		fmt.Println("failed to remove thumbnail")
	}

	err = db.p.Del(m.Hash)
	if err != nil {
		fmt.Println("failed to remove file")
	}

	err = db.m.RemoveThumbnails(m)
	if err != nil {
		fmt.Println("failed to delete thumbnail records")
		return err
	}

	return db.m.DeleteMetaById(m.ID)
}

func (db DB) validate() error {
	if db.m == nil {
		return errors.New("no storage specified")
	}

	if db.p == nil {
		return errors.New("no persistence specified")
	}

	return nil
}

func validateMeta(meta *File) error {
	if meta.Hash == "" {
		return errors.New("no hash specified")
	}

	if meta.Size == 0 {
		return errors.New("no size specified")
	}

	if meta.Name == "" {
		return errors.New("no name specified")
	}

	return nil
}

func (db DB) GetData(h string) (io.ReadCloser, error) {
	return db.p.Get(h)
}

func (db DB) finish(m *File) error {
	w, err := db.p.Get(m.Hash)
	if err != nil {
		return err
	}
	defer w.Close()

	h, err := calcHash(w)
	if err != nil {
		return err
	}

	if h != m.Hash {
		fmt.Println("got:", m.Hash, "expected", h)
		return errors.New("hash does not match")
	}

	err = db.process(m)
	if err != nil {
		return nil
	}

	return nil
}

func (db DB) process(m *File) error {
	fmt.Println("processing")
	for i, p := range db.px {
		fmt.Println("processing", i)
		err := p.Process(db, m)
		if err != nil {
			fmt.Println("processing error", err)
		}
	}
	return nil
}

func (db DB) Path(s string) string {
	return db.p.Path(s)
}

func (db DB) store(m *File, rc io.ReadCloser) error {
	fmt.Println("storing")
	if m.Finished() {
		return errors.New("file already uploaded")
	}

	writer, err := db.p.Put(m.Hash)
	if err != nil {
		return err
	}

	n, _ := io.Copy(writer, rc)
	m.BytesReceived += int(n)

	fmt.Println("after br", m.BytesReceived, m.Size)
	if m.Finished() {
		fmt.Println("setting slug")
		m.Slug = randStr(5)
	} else {

		fmt.Println("NOT setting slug")
	}

	if err := db.m.StoreMeta(m); err != nil {
		return err
	}

	if !m.Finished() {
		return errors.New("got partial data")
	}

	return nil
}

func calcHash(src io.Reader) (string, error) {
	h := sha1.New()

	_, err := io.Copy(h, src)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

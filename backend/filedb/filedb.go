package filedb

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type persister interface {
	Del(string) error
	Put(string) (io.WriteCloser, error)
	Get(string) (io.ReadCloser, error)
}

type metaStorer interface {
	FetchMetaWithSlug(string) (*Meta, error)
	FetchMeta(string) (*Meta, error)
	StoreMeta(*Meta) error
	ListPage(int) ([]*Meta, error)
}

// FileDB implements a upload server.
type FileDB struct {
	p persister
	m metaStorer
}

// NewFileDB returns an instance of a FileDB.
func NewFileDB(p persister, m metaStorer) FileDB {
	return FileDB{
		m: m,
		p: p,
	}
}

func (db FileDB) List(page int) ([]*Meta, error) {
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

func (db FileDB) Write(hash string, rc io.ReadCloser) (*Meta, error) {
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

func (db FileDB) Read(hash string, wc io.Writer) error {
	if err := db.validate(); err != nil {
		return err
	}

	m, err := db.m.FetchMeta(hash)
	if err != nil {
		return err
	}

	if !m.finished() {
		return errors.New("file incomplete")
	}

	reader, err := db.p.Get(hash)

	if err != nil {
		return err
	}

	_, err = io.Copy(wc, reader)

	return err
}

func (db FileDB) StoreMeta(meta Meta) error {
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

func (db FileDB) FetchMeta(h string) (*Meta, error) {
	return db.fetch(h)
}

func (db FileDB) fetch(hash string) (*Meta, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	if hash == "" {
		return nil, errors.New("no hash specified")
	}

	return db.m.FetchMeta(hash)
}

func (db FileDB) FetchMetaWithSlug(slug string) (*Meta, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	if slug == "" {
		return nil, errors.New("no slug specified")
	}

	return db.m.FetchMetaWithSlug(slug)
}

func (db FileDB) validate() error {
	if db.m == nil {
		return errors.New("no storage specified")
	}

	if db.p == nil {
		return errors.New("no persistence specified")
	}

	return nil
}

func validateMeta(meta *Meta) error {
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

func (db FileDB) finish(m *Meta) error {
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
		return errors.New("hash does not match")
	}

	err = db.process(m)
	if err != nil {
		// fmt.Println("error processing:", err.Error())
		return nil
	}

	return nil
}

type bufSeeker struct {
	*bytes.Buffer
}

func (_ bufSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (db FileDB) process(m *Meta) error {
	return nil
}

func (db FileDB) store(m *Meta, rc io.ReadCloser) error {
	if m.finished() {
		return errors.New("file already uploaded")
	}

	writer, err := db.p.Put(m.Hash)
	if err != nil {
		return err
	}

	n, _ := io.Copy(writer, rc)
	m.BytesReceived += int(n)

	if m.finished() {
		m.Slug = randStr(5)
	}

	if err := db.m.StoreMeta(m); err != nil {
		return err
	}

	if !m.finished() {
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

package filedb

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

// Persister reads and writes data.
type Persister interface {
	Put(string) (io.WriteCloser, error)
	Get(string) (io.ReadCloser, error)
}

// Metastorer reads and writes Meta data.
type MetaStorer interface {
	FetchMeta(string) (*Meta, error)
	StoreMeta(Meta) error
}

// FileDB implements a upload server.
type FileDB struct {
	p Persister
	m MetaStorer
}

// NewFileDB returns an instance of a FileDB.
func NewFileDB(p Persister, m MetaStorer) FileDB {
	return FileDB{
		m: m,
		p: p,
	}
}

func (db FileDB) Write(hash string, rc io.ReadCloser) (*Meta, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	m, err := db.FetchMeta(hash)
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

	return db.m.StoreMeta(meta)
}

func (db FileDB) FetchMeta(hash string) (*Meta, error) {
	if err := db.validate(); err != nil {
		return nil, err
	}

	if hash == "" {
		return nil, errors.New("no hash specified")
	}

	return db.m.FetchMeta(hash)
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

	if err := db.m.StoreMeta(*m); err != nil {
		return err
	}

	if !m.finished() {
		return errors.New("got partial data")
	}

	return nil
}

func calcHash(src io.Reader) (string, error) {
	h := sha1.New()

	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

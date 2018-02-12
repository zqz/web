package filedb

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Persister interface {
	Put(string) (io.WriteCloser, error)
	Get(string) (io.ReadCloser, error)
}

type MetaStorer interface {
	FetchMeta(string) (*Meta, error)
	StoreMeta(Meta) error
}

type FileDB struct {
	p Persister
	m MetaStorer
}

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

	if m.Size == m.BytesReceived {
		return m, errors.New("file already uploaded")
	}

	writer, err := db.p.Put(hash)

	if err != nil {
		return nil, err
	}

	n, _ := io.Copy(writer, rc)

	m.BytesReceived += int(n)

	if m.finished() {
		m.Slug = randStr(5)
	}

	err = db.m.StoreMeta(*m)
	if err != nil {
		return nil, err
	}

	if !m.finished() {
		return m, errors.New("got partial data")
	}

	w, err := db.p.Get(m.Hash)
	if err != nil {
		return nil, err
	}
	defer w.Close()

	h, err := calcHash(w)

	if h != m.Hash {
		return nil, errors.New("hash does not match")
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

func calcHash(src io.Reader) (string, error) {
	h := sha1.New()

	if _, err := io.Copy(h, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

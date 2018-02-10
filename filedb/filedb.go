package filedb

import (
	"errors"
	"io"
)

type Persister interface {
	Put(string) (io.WriteCloser, error)
	Get(string) (io.ReadCloser, error)
}

type MetaStorer interface {
	FetchMeta(string) (*Meta, error)
	StoreMeta(*Meta) error
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

func (db FileDB) Write(hash string, rc io.ReadCloser) (int64, error) {
	if db.p == nil {
		return 0, errors.New("no persistence specified")
	}

	writer, err := db.p.Put(hash)

	if err != nil {
		return 0, nil
	}

	return io.Copy(writer, rc)
}

func (db FileDB) Read(hash string, wc io.Writer) (int64, error) {
	if db.p == nil {
		return 0, errors.New("no persistence specified")
	}

	reader, err := db.p.Get(hash)

	if err != nil {
		return 0, err
	}

	return io.Copy(wc, reader)
}

func (db FileDB) StoreMeta(meta *Meta) error {
	if db.m == nil {
		return errors.New("no storage specified")
	}

	if meta.Hash == "" {
		return errors.New("no hash specified")
	}

	if meta.Size == 0 {
		return errors.New("no size specified")
	}

	if meta.Name == "" {
		return errors.New("no name specified")
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
	if db.m == nil {
		return nil, errors.New("no storage specified")
	}

	return db.m.FetchMeta(hash)
}

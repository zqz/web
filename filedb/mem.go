package filedb

import (
	"bytes"
	"errors"
	"io"
)

type nopReadCloser struct {
	io.Reader
}

type nopWriteCloser struct {
	io.Writer
}

func (nopReadCloser) Close() error  { return nil }
func (nopWriteCloser) Close() error { return nil }

type MemoryPersistence struct {
	entries map[string]*bytes.Buffer
}

type MemoryMetaStorage struct {
	entries map[string]*Meta
}

func NewMemoryMetaStorage() MemoryMetaStorage {
	return MemoryMetaStorage{
		entries: make(map[string]*Meta),
	}
}

func NewMemoryPersistence() MemoryPersistence {
	e := make(map[string]*bytes.Buffer)

	return MemoryPersistence{
		entries: e,
	}
}

func (m MemoryPersistence) Put(hash string) (io.WriteCloser, error) {
	b, ok := m.entries[hash]

	if !ok {
		b = new(bytes.Buffer)
		m.entries[hash] = b
	}

	return nopWriteCloser{b}, nil
}

func (m MemoryPersistence) Get(hash string) (io.ReadCloser, error) {
	b, ok := m.entries[hash]

	if !ok {
		return nil, errors.New("no file with hash: " + hash)
	}

	return nopReadCloser{b}, nil
}

func (m MemoryMetaStorage) FetchMeta(hash string) (*Meta, error) {
	meta, ok := m.entries[hash]

	if !ok {
		return nil, errors.New("file not found")
	}

	return meta, nil
}

func (m MemoryMetaStorage) StoreMeta(meta Meta) error {
	// m.entries[meta.Hash]

	// 	if ok {
	// 		return errors.New("file already exists")
	// 	}

	m.entries[meta.Hash] = &meta

	return nil
}

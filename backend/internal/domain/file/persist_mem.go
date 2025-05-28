package file

import (
	"bytes"
	"errors"
	"io"

	"github.com/davecgh/go-spew/spew"
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

func NewMemoryPersistence() MemoryPersistence {
	e := make(map[string]*bytes.Buffer)

	return MemoryPersistence{
		entries: e,
	}
}

func (m MemoryPersistence) Del(hash string) error {
	if _, ok := m.entries[hash]; ok {
		return m.Del(hash)
	}
	return nil
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

	spew.Dump(m.entries)

	if !ok {
		return nil, errors.New("no file with hash: " + hash)
	}

	data := b.Bytes()
	buf := bytes.NewBuffer(data)

	return nopReadCloser{buf}, nil
}

func (m MemoryPersistence) Path(hash string) string {
	return hash
}

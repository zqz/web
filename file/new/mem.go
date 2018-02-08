package new

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

type MemoryPersistance struct {
	entries map[string]*bytes.Buffer
}

func NewMemoryPersistance() MemoryPersistance {
	e := make(map[string]*bytes.Buffer)

	return MemoryPersistance{
		entries: e,
	}
}

func (m MemoryPersistance) Put(hash string) (io.WriteCloser, error) {
	b, ok := m.entries[hash]

	if !ok {
		b = new(bytes.Buffer)
		m.entries[hash] = b
	}

	return nopWriteCloser{b}, nil
}

func (m MemoryPersistance) Get(hash string) (io.ReadCloser, error) {
	b, ok := m.entries[hash]

	if !ok {
		return nil, errors.New("no file with hash: " + hash)
	}

	return nopReadCloser{b}, nil
}

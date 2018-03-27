package filedb

import (
	"io"
	"os"
)

type DiskPersistence struct {
}

var fileFlags int = os.O_APPEND | os.O_CREATE | os.O_WRONLY
var fileMode os.FileMode = 0600

func init() {
	os.MkdirAll(path(""), 0744)
}

func path(hash string) string {
	return "/tmp/zqz/" + hash
}

func NewDiskPersistence() DiskPersistence {
	return DiskPersistence{}
}

func (DiskPersistence) Put(hash string) (io.WriteCloser, error) {
	p := path(hash)
	return os.OpenFile(p, fileFlags, fileMode)
}

func (DiskPersistence) Get(hash string) (io.ReadCloser, error) {
	p := path(hash)
	return os.Open(p)
}

// func (m MemoryMetaStorer) FetchMeta(hash string) (*Meta, error) {
// 	meta, ok := m.entries[hash]

// 	if !ok {
// 		return nil, errors.New("file not found")
// 	}

// 	return meta, nil
// }

// func (m MemoryMetaStorer) StoreMeta(meta Meta) error {
// 	_, ok := m.entries[meta.Hash]

// 	if ok {
// 		return errors.New("file already exists")
// 	}

// 	m.entries[meta.Hash] = &meta

// 	return nil
// }

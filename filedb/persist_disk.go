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
	return "./files/" + hash
}

func NewDiskPersistence() DiskPersistence {
	return DiskPersistence{}
}

func (DiskPersistence) Del(hash string) error {
	p := path(hash)
	return os.Remove(p)
}

func (DiskPersistence) Put(hash string) (io.WriteCloser, error) {
	p := path(hash)
	return os.OpenFile(p, fileFlags, fileMode)
}

func (DiskPersistence) Get(hash string) (io.ReadCloser, error) {
	p := path(hash)
	return os.Open(p)
}

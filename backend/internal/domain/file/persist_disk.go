package filedb

import (
	"io"
	"os"
	"path/filepath"
)

type DiskPersistence struct {
	DstPath string
}

var fileFlags int = os.O_APPEND | os.O_CREATE | os.O_WRONLY
var fileMode os.FileMode = 0600
var dirMode os.FileMode = 0744

func NewDiskPersistence(path string) (DiskPersistence, error) {
	d := DiskPersistence{DstPath: path}

	if ok, err := d.init(); !ok {
		return d, err
	}

	return d, nil
}

func (d DiskPersistence) Del(hash string) error {
	p := d.path(hash)
	return os.Remove(p)
}

func (d DiskPersistence) Put(hash string) (io.WriteCloser, error) {
	p := d.path(hash)

	return os.OpenFile(p, fileFlags, fileMode)
}

func (d DiskPersistence) Get(hash string) (io.ReadCloser, error) {
	p := d.path(hash)
	return os.Open(p)
}

func (d DiskPersistence) init() (bool, error) {
	err := os.MkdirAll(d.path(""), dirMode)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (d DiskPersistence) Path(hash string) string {
	return d.path(hash)
}

func (d DiskPersistence) path(hash string) string {
	return filepath.Join(d.DstPath, hash)
}

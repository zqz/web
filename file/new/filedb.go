package new

import (
	"fmt"
	"io"

	"github.com/davecgh/go-spew/spew"
)

type Persister interface {
	Put(string) (io.WriteCloser, error)
	Get(string) (io.ReadCloser, error)
}

type MetaStorer interface {
	FetchMeta(string) (*Meta, error)
	StoreMeta(Meta) (int, error)
}

type FileDB struct {
	p Persister
	m MetaStorer
}

func (db FileDB) Write(hash string, rc io.ReadCloser) error {
	writer, err := db.p.Put(hash)

	if err != nil {
		return nil
	}

	n, err := io.Copy(writer, rc)
	fmt.Println("wrote", n, "bytes")

	spew.Dump(db.p)

	//writer.Close()

	return err
}

func (db FileDB) Read(hash string, wc io.WriteCloser) error {
	reader, err := db.p.Get(hash)

	if err != nil {
		return err
	}

	n, err := io.Copy(wc, reader)
	fmt.Println("read", n, "bytes")

	//	reader.Close()

	return err
}

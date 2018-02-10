package filedb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteNoPersistence(t *testing.T) {
	db := FileDB{}
	n, err := db.Write("hash", nopReadCloser{})
	assert.Equal(t, "no persistence specified", err.Error())
	assert.Equal(t, int64(0), n)
}

func TestReadNoPersistence(t *testing.T) {
	db := FileDB{}
	n, err := db.Read("hash", nopWriteCloser{})
	assert.Equal(t, "no persistence specified", err.Error())
	assert.Equal(t, int64(0), n)
}

func TestHelloWorld(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
	}

	str := "hello world"
	rc := nopReadCloser{bytes.NewBufferString(str)}

	n, err := db.Write("hash", rc)
	assert.Nil(t, err)
	assert.Equal(t, len(str), int(n))

	var b bytes.Buffer
	wc := nopWriteCloser{&b}

	n, err = db.Read("hash", wc)
	assert.Nil(t, err)
	assert.Equal(t, len(str), int(n))
	assert.Equal(t, "hello world", b.String())
}

func TestStoreMetaNoStorage(t *testing.T) {
	db := FileDB{}
	m := Meta{}

	err := db.StoreMeta(m)

	assert.Equal(t, "no storage specified", err.Error())
}

func TestStoreMetaNoHash(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}
	m := Meta{}

	err := db.StoreMeta(m)
	assert.Equal(t, "no hash specified", err.Error())
}

func TestStoreMetaNoSize(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	m := Meta{
		Hash: "foo",
	}

	err := db.StoreMeta(m)
	assert.Equal(t, "no size specified", err.Error())
}

func TestStoreMetaNoName(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	m := Meta{
		Hash: "foo",
		Size: 123,
	}

	err := db.StoreMeta(m)
	assert.Equal(t, "no name specified", err.Error())
}

func TestStoreMeta(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	m := Meta{
		Hash: "foo",
		Size: 123,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	testMeta, err := db.FetchMeta("foo")

	assert.Equal(t, &m, testMeta)
}

func TestStoreMetaCantChangeSize(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	m := Meta{
		Hash: "foo",
		Size: 123,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	m.Size = 321

	err = db.StoreMeta(m)

	assert.Equal(t, "can not change file size", err.Error())
}

func TestFetchMetaNoStorage(t *testing.T) {
	db := FileDB{}

	meta, err := db.FetchMeta("hash")

	assert.Equal(t, "no storage specified", err.Error())
	assert.Nil(t, meta)
}

func TestFetchMetaNoHash(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	meta, err := db.FetchMeta("")

	assert.Equal(t, "no hash specified", err.Error())
	assert.Nil(t, meta)
}

func TestFetchMetaNoMeta(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	meta, err := db.FetchMeta("foo")

	assert.Equal(t, "file not found", err.Error())
	assert.Nil(t, meta)
}

func TestFetchMeta(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	m := Meta{
		Hash: "foo",
		Size: 123,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	meta, err := db.FetchMeta("foo")
	assert.Nil(t, err)

	assert.Equal(t, meta.Hash, m.Hash)
	assert.Equal(t, meta.Size, m.Size)
	assert.Equal(t, meta.Name, m.Name)
	assert.Equal(t, 0, meta.BytesReceived)
}

func TestCantWriteOnceReceivedAllData(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	m := Meta{
		Hash: "foo",
		Size: 5,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	n, err := db.Write("hash", rc)

	assert.Nil(t, err)
	assert.Equal(t, 5, int(n))

	meta, err := db.FetchMeta(m.Hash)

	assert.Equal(t, 5, meta.BytesReceived)
}

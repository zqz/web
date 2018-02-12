package filedb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteNoPersistence(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	m, err := db.Write("hash", nopReadCloser{})
	assert.Nil(t, m)
	assert.Equal(t, "no persistence specified", err.Error())
}

func TestReadNoPersistence(t *testing.T) {
	db := FileDB{
		m: NewMemoryMetaStorage(),
	}

	err := db.Read("hash", nopWriteCloser{})
	assert.Equal(t, "no persistence specified", err.Error())
}

func TestStoreMetaNoStorage(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
	}
	m := Meta{}

	err := db.StoreMeta(m)

	assert.Equal(t, "no storage specified", err.Error())
}

func TestStoreMetaNoHash(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}
	m := Meta{}

	err := db.StoreMeta(m)
	assert.Equal(t, "no hash specified", err.Error())
}

func TestStoreMetaNoSize(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
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
		p: NewMemoryPersistence(),
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
		p: NewMemoryPersistence(),
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
		p: NewMemoryPersistence(),
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
	db := FileDB{
		p: NewMemoryPersistence(),
	}

	meta, err := db.FetchMeta("hash")

	assert.Equal(t, "no storage specified", err.Error())
	assert.Nil(t, meta)
}

func TestFetchMetaNoHash(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	meta, err := db.FetchMeta("")

	assert.Equal(t, "no hash specified", err.Error())
	assert.Nil(t, meta)
}

func TestFetchMetaNoMeta(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	meta, err := db.FetchMeta("foo")

	assert.Equal(t, "file not found", err.Error())
	assert.Nil(t, meta)
}

func TestFetchMeta(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
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

func TestWriteSuccess(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Hash: hash,
		Size: 5,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	meta, err := db.Write(hash, rc)

	assert.Nil(t, err)
	assert.Equal(t, 5, meta.Size)
	assert.Equal(t, 5, meta.BytesReceived)
	assert.Equal(t, "foobar", meta.Name)
	assert.Equal(t, hash, meta.Hash)
}

func TestReturnErrorOnBadHash(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "badhash"

	m := Meta{
		Hash: hash,
		Size: 5,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	meta, err := db.Write(hash, rc)

	assert.Nil(t, meta)
	assert.Equal(t, "hash does not match", err.Error())
}

func TestCanNotWriteOnceReceivedAllData(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Hash: hash,
		Size: 5,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	meta, err := db.Write(hash, rc)

	assert.Nil(t, err)
	assert.Equal(t, 5, meta.BytesReceived)

	rc = nopReadCloser{bytes.NewBufferString(str)}
	meta, err = db.Write(hash, rc)

	assert.Equal(t, "file already uploaded", err.Error())
	assert.Equal(t, 5, meta.BytesReceived)
}

func TestReadPartial(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Hash: hash,
		Size: 5,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	str := "byt"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	_, err = db.Write(hash, rc)
	assert.NotNil(t, err) // partial upload

	var b bytes.Buffer
	wc := nopWriteCloser{&b}
	err = db.Read(hash, wc)

	assert.Equal(t, "file incomplete", err.Error())
	assert.Equal(t, "", b.String())
}

func TestFull(t *testing.T) {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Hash: hash,
		Size: 5,
		Name: "foobar",
	}

	err := db.StoreMeta(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	meta, err := db.Write(hash, rc)

	assert.Nil(t, err)
	assert.Equal(t, 5, meta.Size)
	assert.Equal(t, 5, meta.BytesReceived)
	assert.Equal(t, "foobar", meta.Name)
	assert.Equal(t, hash, meta.Hash)

	var b bytes.Buffer
	wc := nopWriteCloser{&b}
	err = db.Read(hash, wc)

	assert.Nil(t, err)
	assert.Equal(t, "bytes", b.String())
}

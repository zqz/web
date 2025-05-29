package file

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriteNoPersistence(t *testing.T) {
	db := DB{
		m: NewMemoryMetaStorage(),
	}

	m, err := db.Write("hash", nopReadCloser{})
	assert.Nil(t, m)
	assert.Equal(t, "no persistence specified", err.Error())
}

func TestReadNoPersistence(t *testing.T) {
	db := DB{
		m: NewMemoryMetaStorage(),
	}

	err := db.Read("hash", nopWriteCloser{})
	assert.Equal(t, "no persistence specified", err.Error())
}

func TestCreateNoStorage(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
	}
	m := File{}

	err := db.Create(m)

	assert.Equal(t, "no storage specified", err.Error())
}

func TestCreateNoHash(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}
	m := File{}

	err := db.Create(m)
	assert.Equal(t, "no hash specified", err.Error())
}

func TestCreateNoSize(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	m := File{}
	m.Hash = "foo"

	err := db.Create(m)
	assert.Equal(t, "no size specified", err.Error())
}

func TestCreateNoName(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	m := File{}
	m.Hash = "foo"
	m.Size = 123

	err := db.Create(m)
	assert.Equal(t, "no name specified", err.Error())
}

func TestCreate(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	m := File{}
	m.Hash = "foo"
	m.Size = 123
	m.Name = "foobar"

	err := db.Create(m)
	assert.Nil(t, err)

	testMeta, err := db.FetchByHash("foo")
	m.ID = testMeta.ID // not sure what id would be in the test

	assert.Equal(t, &m, testMeta)
}

func TestCreateCantChangeSize(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	m := File{}
	m.Hash = "foo"
	m.Size = 123
	m.Name = "foobar"

	err := db.Create(m)
	assert.Nil(t, err)

	m.Size = 321

	err = db.Create(m)

	assert.Equal(t, "can not change file size", err.Error())
}

func TestFetchByHashNoStorage(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
	}

	meta, err := db.FetchByHash("hash")

	assert.Equal(t, "no storage specified", err.Error())
	assert.Nil(t, meta)
}

func TestFetchByHashNoHash(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	meta, err := db.FetchByHash("")

	assert.Equal(t, "no hash specified", err.Error())
	assert.Nil(t, meta)
}

func TestFetchByHashNoMeta(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	meta, err := db.FetchByHash("foo")

	assert.Equal(t, "file not found", err.Error())
	assert.Nil(t, meta)
}

func TestFetchByHash(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	m := File{}
	m.Hash = "foo"
	m.Size = 123
	m.Name = "foobar"
	err := db.Create(m)
	assert.Nil(t, err)

	meta, err := db.FetchByHash("foo")
	assert.Nil(t, err)

	assert.Equal(t, meta.Hash, m.Hash)
	assert.Equal(t, meta.Size, m.Size)
	assert.Equal(t, meta.Name, m.Name)
	assert.Equal(t, 0, meta.BytesReceived)
}

func TestWriteSuccess(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
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

func TestFetchBySlug(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"
	m.Slug = "doo"
	err := db.Create(m)
	assert.Nil(t, err)

	meta, err := db.FetchBySlug("doo")

	assert.Nil(t, err)
	assert.Equal(t, 5, meta.Size)
	assert.Equal(t, "foobar", meta.Name)
	assert.Equal(t, hash, meta.Hash)
}

func TestReturnErrorOnBadHash(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "badhash"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	meta, err := db.Write(hash, rc)

	assert.Nil(t, meta)
	assert.Equal(t, "hash does not match", err.Error())
}

func TestCanNotWriteOnceReceivedAllData(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
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
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
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
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
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

func TestListPartial(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
	assert.Nil(t, err)

	str := "byt"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	db.Write(hash, rc)

	metas, err := db.List(0)

	assert.Equal(t, 0, len(metas))
}

func TestList(t *testing.T) {
	db := DB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := File{}
	m.Hash = hash
	m.Size = 5
	m.Name = "foobar"

	err := db.Create(m)
	assert.Nil(t, err)

	str := "bytes"
	rc := nopReadCloser{bytes.NewBufferString(str)}
	db.Write(hash, rc)

	metas, err := db.List(0)

	assert.Equal(t, 1, len(metas))
}

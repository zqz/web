package filedb

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testServer() *httptest.Server {
	db := FileDB{
		p: NewMemoryPersistence(),
		m: NewMemoryMetaStorage(),
	}

	s := Server{
		db: db,
	}

	ts := httptest.NewServer(s.Router())

	return ts
}

func get(ts *httptest.Server, path string) *http.Response {
	res, _ := http.Get(ts.URL + path)

	return res
}

func post(ts *httptest.Server, path string, o interface{}) *http.Response {
	b, _ := json.Marshal(&o)

	buf := bytes.NewBuffer(b)

	res, _ := http.Post(ts.URL+path, "application/json", buf)

	return res
}

func TestPostMeta(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Name: "foo",
		Size: 5,
		Hash: hash,
	}

	res := post(ts, "/meta", m)

	assert.Equal(t, 200, res.StatusCode)

	responseBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	b, _ := json.Marshal(&m)
	assert.Equal(t, string(b), string(responseBody))
}

func TestGetMetaNotFound(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	res := get(ts, "/meta/"+hash)

	assert.Equal(t, 404, res.StatusCode)
}

func TestGetMetaFound(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Name: "foo",
		Size: 5,
		Hash: hash,
	}

	post(ts, "/meta", m)

	res := get(ts, "/meta/"+hash)

	assert.Equal(t, 200, res.StatusCode)
	responseBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	b, _ := json.Marshal(&m)
	assert.Equal(t, string(b), string(responseBody))
}

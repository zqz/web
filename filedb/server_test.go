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

func toJSON(o interface{}) string {
	b, _ := json.Marshal(&o)

	return string(b)
}

func errorMessage(res *http.Response) string {
	o := struct {
		Message string `json:"message"`
	}{}

	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	json.Unmarshal(b, &o)

	return o.Message
}

func readBody(res *http.Response) string {
	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	return string(b)
}

func get(ts *httptest.Server, path string) *http.Response {
	res, _ := http.Get(ts.URL + path)

	return res
}

func postData(ts *httptest.Server, path string, data string) *http.Response {
	buf := bytes.NewBufferString(data)
	res, _ := http.Post(ts.URL+path, "application/octet-stream", buf)
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
	assert.Equal(t, "file not found", errorMessage(res))
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
	assert.Equal(t, toJSON(m), readBody(res))
}

func TestPostDataNoMeta(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	res := post(ts, "/data/"+hash, "foobar")

	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, "file not found", errorMessage(res))
}

func TestPostData(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Name: "foo",
		Size: 5,
		Hash: hash,
	}

	post(ts, "/meta", m)

	res := postData(ts, "/data/"+hash, "byt")
	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, "got partial data", errorMessage(res))

	res = postData(ts, "/data/"+hash, "es")
	assert.Equal(t, 200, res.StatusCode)

	json.Unmarshal([]byte(readBody(res)), &m)
	assert.Equal(t, 5, m.BytesReceived)
	assert.NotEmpty(t, m.Slug)
}

func TestGetDataNotFound(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	res := get(ts, "/data/"+hash)
	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, "file not found", errorMessage(res))
}

func TestGetData(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Name:        "foo",
		Size:        5,
		Hash:        hash,
		ContentType: "foo/bar",
	}

	post(ts, "/meta", m)
	postData(ts, "/data/"+hash, "bytes")

	res := get(ts, "/data/"+hash)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "foo/bar", res.Header.Get("Content-Type"))
	assert.Equal(t, hash, res.Header.Get("Etag"))
	assert.Equal(t, "no-cache", res.Header.Get("Cache-Control"))
	assert.Equal(t, "inline; filename="+m.Name, res.Header.Get("Content-Disposition"))
}

func TestGetDataCachedInBrowser(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	m := Meta{
		Name:        "foo",
		Size:        5,
		Hash:        hash,
		ContentType: "foo/bar",
	}

	post(ts, "/meta", m)
	postData(ts, "/data/"+hash, "bytes")

	req, _ := http.NewRequest("GET", ts.URL+"/data/"+hash, nil)
	req.Header.Add("If-None-Match", hash)
	res, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 304, res.StatusCode)
}

func TestGetDataUnknwonCachedInBrowser(t *testing.T) {
	ts := testServer()
	defer ts.Close()

	hash := "daf529a73101c2be626b99fc6938163e7a27620b"

	req, _ := http.NewRequest("GET", ts.URL+"/data/"+hash, nil)
	req.Header.Add("If-None-Match", hash)
	res, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, "file not found", errorMessage(res))
}

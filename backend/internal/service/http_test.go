package service_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/assert"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/service"
)

func TestAuthenticatedAdminEndpoints(t *testing.T) {
	tests := []struct {
		method string
		name   string
		path   string
		status int
	}{
		{"GET", "admin files", "/admin/files", http.StatusOK},
		{"GET", "admin users", "/admin/users", http.StatusOK},
		{"GET", "admin user by id", "/admin/users/1231", http.StatusNotFound},
		{"GET", "admin file by id", "/files/1231", http.StatusNotFound},
		{"GET", "admin edit user by id", "/admin/users/1231/edit", http.StatusNotFound},
		{"GET", "admin edit file by id", "/files/1231/edit", http.StatusNotFound},
	}

	s := service.NewTestServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			response, _ := requestAsAdmin(t, req, &s)

			if response.Code != tt.status {
				t.Errorf("Expected %d, got %d for %s", tt.status, response.Code, tt.path)
			}
		})
	}
}

func TestUnauthenticatedAdminEndpoints(t *testing.T) {
	tests := []struct {
		method string
		name   string
		path   string
	}{
		{"GET", "admin files", "/admin/files"},
		{"GET", "admin users", "/admin/users"},
		{"GET", "admin user by id", "/admin/users/1231"},
		{"GET", "admin file by id", "/admin/files/1231"},
		{"GET", "admin edit user by id", "/admin/users/1231/edit"},
		{"GET", "admin edit file by id", "/admin/files/1231/edit"},
	}

	s := service.NewTestServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			response := requestAsGuest(req, &s)

			if response.Code != http.StatusForbidden {
				t.Errorf("Expected %d, got %d for %s", http.StatusForbidden, response.Code, tt.path)
			}
		})
	}
}

func TestGetRootAsAdmin(t *testing.T) {
	s := service.NewTestServer()
	req, _ := http.NewRequest("GET", "/", nil)
	response, _ := requestAsAuthed(t, req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	assert.NoError(t, err)

	if doc.Find(`[data-user]`).Length() != 1 {
		t.Error("authed user should see the authed info component")
	}
}

func TestGetRootWorks(t *testing.T) {
	s := service.NewTestServer()

	req, _ := http.NewRequest("GET", "/", nil)
	response := requestAsGuest(req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetUnknownFileReturns404(t *testing.T) {
	s := service.NewTestServer()

	req, _ := http.NewRequest("GET", "/files/test123", nil)
	response := requestAsGuest(req, &s)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetFileReturns200(t *testing.T) {
	s := service.NewTestServer()

	f := addFile(t, &s, "test")

	req, _ := http.NewRequest("GET", "/files/"+f.Slug, nil)
	response := requestAsGuest(req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetFilePreviewReturns200(t *testing.T) {
	s := service.NewTestServer()

	f := addFile(t, &s, "test")

	req, _ := http.NewRequest("GET", "/files/"+f.Slug+"/preview", nil)
	response := requestAsGuest(req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetFileDoesNotShowEditForNotAdmin(t *testing.T) {
	s := service.NewTestServer()

	f := addFile(t, &s, "test")

	req, _ := http.NewRequest("GET", "/files/"+f.Slug, nil)
	response := requestAsGuest(req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	assert.NoError(t, err)

	if doc.Find(`[data-edit]`).Length() != 0 {
		t.Error("edit should not be visible for guest")
	}
}

func addFile(t *testing.T, s *service.Server, name string) *file.Meta {
	content := "test"
	hash := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	f := file.Meta{}
	f.ID = 123
	f.Hash = hash
	f.Name = name
	f.Size = len(content)
	err := s.FileDB.StoreMeta(f)
	assert.NoError(t, err)

	// expect to be able to write the data
	_, err = s.FileDB.Write(hash, nopReadCloser{bytes.NewBufferString(content)})
	assert.NoError(t, err)

	// expect to get the data back
	_, err = s.FileDB.GetData(hash)
	assert.NoError(t, err)

	// expect to be able to get the meta back
	f2, err := s.FileDB.FetchMeta(hash)
	assert.NoError(t, err)

	return f2
}

func requestAsAdmin(t *testing.T, req *http.Request, s *service.Server) (*httptest.ResponseRecorder, *user.User) {
	u := user.User{}
	u.Name = "test"
	u.Email = "qdylanj@gmail.com"
	u.Provider = "goggle"
	u.ProviderID = "123"

	return requestAsUser(t, req, s, &u)
}

func requestAsAuthed(t *testing.T, req *http.Request, s *service.Server) (*httptest.ResponseRecorder, *user.User) {
	u := user.User{}
	u.Name = "test"
	u.Email = "test@site.com"
	u.Provider = "goggle"
	u.ProviderID = "123"

	return requestAsUser(t, req, s, &u)
}

func requestAsUser(t *testing.T, req *http.Request, s *service.Server, u *user.User) (*httptest.ResponseRecorder, *user.User) {
	err := s.UserDB.Create(u)
	assert.NoError(t, err)

	u2, err := s.UserDB.FindByProviderId("123")
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	gothic.StoreInSession("user_id", strconv.Itoa(u2.ID), req, w)

	for _, cookie := range w.Result().Cookies() {
		req.AddCookie(cookie)
	}

	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr, u2
}

func requestAsGuest(req *http.Request, s *service.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

type nopReadCloser struct {
	io.Reader
}

type nopWriteCloser struct {
	io.Writer
}

func (nopReadCloser) Close() error  { return nil }
func (nopWriteCloser) Close() error { return nil }

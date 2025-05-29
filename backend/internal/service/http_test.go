package service_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/assert"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/service"
)

func TestEdgecaseEndpoints(t *testing.T) {
	s := service.NewTestServer()

	tests := []struct {
		method string
		name   string
		path   string
		status int
	}{
		{"GET", "non existing route", "/foo/bar", http.StatusNotFound},
		{"POST", "non existing method existing route", "/", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			response, _ := requestAsAdmin(t, req, &s)

			if response.Code != tt.status {
				t.Errorf("Expected %d, got %d for %s", tt.status, response.Code, tt.path)
				fmt.Println(string(response.Body.String()))
			}
		})
	}
}

func TestAuthenticatedAdminEndpoints(t *testing.T) {
	s := service.NewTestServer()
	f := addFile(t, &s, "test")
	u := addUser(t, &s, "test")

	userId := strconv.Itoa(u.ID)

	tests := []struct {
		method string
		name   string
		path   string
		status int
	}{
		{"GET", "admin files", "/admin/files", http.StatusOK},
		{"GET", "admin users", "/admin/users", http.StatusOK},

		{"GET", "admin user by id", "/admin/users/" + userId, http.StatusOK},
		{"GET", "admin user edit by id", "/admin/users/" + userId + "/edit", http.StatusOK},
		{"GET", "admin file by id", "/files/" + f.Slug, http.StatusOK},
		{"GET", "admin file edit by id", "/admin/files/" + f.Slug + "/edit", http.StatusOK},

		{"GET", "admin unk user by id", "/admin/users/1231", http.StatusNotFound},
		{"GET", "admin unk user edit by id", "/admin/users/1231/edit", http.StatusNotFound},
		{"GET", "admin unk file by id", "/files/1231", http.StatusNotFound},
		{"GET", "admin unk file edit by id", "/admin/files/1231/edit", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			response, _ := requestAsAdmin(t, req, &s)

			if response.Code != tt.status {
				t.Errorf("Expected %d, got %d for %s", tt.status, response.Code, tt.path)
				fmt.Println(string(response.Body.String()))
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

func TestAdminViewEditFile(t *testing.T) {
	s := service.NewTestServer()
	f := addFile(t, &s, "test")

	req, _ := http.NewRequest("GET", "/admin/files/"+f.Slug+"/edit", nil)
	response, _ := requestAsAdmin(t, req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	assert.NoError(t, err)

	name, exists := doc.Find(`input[id="form_name"]`).Attr("value")
	if !exists {
		t.Error("expected to find input with name on the page")
	}

	if name != f.Name {
		t.Errorf("file name should be accessable, got: '%s', expected: '%s'", name, f.Name)
	}
}

func TestAdminEditFileWorks(t *testing.T) {
	s := service.NewTestServer()
	f := addFile(t, &s, "test")

	newFileName := "new-name"
	newSlug := "new-slug"
	newComment := "new-comment"
	newPrivate := "true"

	formData := url.Values{}
	formData.Set("name", newFileName)
	formData.Set("slug", newSlug)
	formData.Set("comment", newComment)
	formData.Set("private", newPrivate)

	// Create request with form data
	data := strings.NewReader(formData.Encode())
	req, _ := http.NewRequest("POST", "/admin/files/"+f.Slug, data)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, _ := requestAsAdmin(t, req, &s)
	checkResponseCode(t, http.StatusFound, response.Code)

	req, _ = http.NewRequest("GET", "/files/"+newSlug, nil)
	response, _ = requestAsAdmin(t, req, &s)
	checkResponseCode(t, http.StatusOK, response.Code)
	doc, err := goquery.NewDocumentFromReader(response.Body)
	assert.NoError(t, err)

	name := doc.Find(`span[data-name]`).Text()
	slug := doc.Find(`span[data-slug]`).Text()
	comment := doc.Find(`span[data-comment]`).Text()
	private := doc.Find(`span[data-private]`).Text()

	if name != newFileName {
		t.Errorf("file name should have been updated, got: '%s', expected: '%s'", name, newFileName)
	}

	if slug != newSlug {
		t.Errorf("file slug should have been updated, got: '%s', expected: '%s'", slug, newSlug)
	}

	if comment != newComment {
		t.Errorf("file comment should have been updated, got: '%s', expected: '%s'", comment, newComment)
	}

	if private != newPrivate {
		t.Errorf("file private should have been updated, got: '%s', expected: '%s'", private, newPrivate)
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

func addUser(t *testing.T, s *service.Server, name string) *user.User {
	u := user.User{}
	u.Name = "test"
	u.Email = "test@site.com"
	u.Provider = "goggle"
	u.ProviderID = "xa1123"

	err := s.UserDB.Create(&u)
	assert.NoError(t, err)

	return &u
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
	u.ProviderID = "321123"

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

	u2, err := s.UserDB.FindByProviderId(u.ProviderID)
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

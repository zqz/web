package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/shared/helper"
	"github.com/zqz/web/backend/render"
)

// Server implements an HTTP File Uploading Server.
type Server struct {
	db  file.FileDB
	udb user.DB
}

// NewServer returns a new Server.
func NewServer(db file.FileDB, udb user.DB) Server {
	return Server{
		db:  db,
		udb: udb,
	}
}

func (s Server) postMeta(w http.ResponseWriter, r *http.Request) {
	var m *file.Meta
	var err error

	u := helper.GetUserFromContext(r.Context())

	if m, err = parseMeta(r.Body); err != nil {
		render.Error(w, "failed to read request")
		return
	}

	spew.Dump(u)

	if m2, err := s.db.FetchMeta(m.Hash); err == nil {
		render.JSON(w, m2)
		return
	}

	if u != nil {
		m.UserID = u.ID
	}

	if err = s.db.StoreMeta(*m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.Error(w, err.Error())
		return
	}

	render.JSON(w, m)
}

func (s Server) getMeta(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	render.JSON(w, meta)
}

func (s Server) postData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	if meta.Finished() {
		w.WriteHeader(http.StatusConflict)
		render.Error(w, "file already uploaded")
		return
	}

	_, err = s.db.Write(hash, r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		render.Error(w, err.Error())
		return
	}

	render.JSON(w, meta)
}

func queryParamInt(r *http.Request, key string) (int, error) {
	v := r.URL.Query().Get(key)
	return strconv.Atoi(v)
}

func (s Server) files(w http.ResponseWriter, r *http.Request) {
	page, err := queryParamInt(r, "page")
	if err != nil {
		page = 0
	}

	m, err := s.db.List(page)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("err", err.Error())
		render.Error(w, "error loading file list")
		return
	}

	render.JSON(w, m)
}

func (s Server) sendfile(hash string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Etag", hash)
	w.Header().Set("Cache-Control", "no-cache")
	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, hash) {
			// go lib.TrackDownload(f.DB, file.ID, r, true)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	rr, err := s.db.GetData(hash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("err:", err.Error())
		return
	}

	io.Copy(w, rr)
}

func (s Server) downloadThumbnail(meta *file.Meta, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename=thumb_"+meta.Name)
	s.sendfile(meta.Thumbnail, w, r)
}

func (s Server) downloadFile(meta *file.Meta, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename="+meta.Name)
	s.sendfile(meta.Hash, w, r)
}

func (s Server) getThumbnailDataWithSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	meta, err := s.db.FetchMetaWithSlug(slug)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	s.downloadThumbnail(meta, w, r)
}

func (s Server) getDataWithSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	meta, err := s.db.FetchMetaWithSlug(slug)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	s.downloadFile(meta, w, r)
}

func (s Server) getData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	s.downloadFile(meta, w, r)
}

func (s Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/files", s.files)
	r.Get("/file/by-hash/{hash}", s.getData)
	r.Get("/file/by-slug/{slug}", s.getDataWithSlug)
	r.Get("/file/by-slug/{slug}/thumbnail", s.getThumbnailDataWithSlug)
	r.Post("/file/{hash}", s.postData)

	r.Get("/meta/by-hash/{hash}", s.getMeta)
	r.Post("/meta", s.postMeta)

	return r
}

func parseMeta(r io.ReadCloser) (*file.Meta, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	m := file.Meta{}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

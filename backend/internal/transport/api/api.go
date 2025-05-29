package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/volatiletech/null/v8"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/shared/helper"
)

// Server implements an HTTP File Uploading Server.
type Server struct {
	db  file.DB
	udb user.DB
}

// NewServer returns a new Server.
func NewServer(db file.DB, udb user.DB) Server {
	return Server{
		db:  db,
		udb: udb,
	}
}

func (s Server) postMeta(w http.ResponseWriter, r *http.Request) {
	var m *file.File
	var err error

	u := helper.GetUserFromContext(r.Context())

	if m, err = parseMeta(r.Body); err != nil {
		Error(w, "failed to read request")
		return
	}

	if m2, err := s.db.FetchByHash(m.Hash); err == nil {
		JSON(w, m2)
		return
	}

	if u != nil {
		m.UserID = null.IntFrom(u.ID)
	}

	if err = s.db.StoreMeta(*m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err.Error())
		return
	}

	JSON(w, m)
}

func (s Server) getMeta(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchByHash(hash)

	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	JSON(w, meta)
}

func (s Server) postData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchByHash(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Error(w, err.Error())
		return
	}

	if meta.Finished() {
		w.WriteHeader(http.StatusConflict)
		Error(w, "file already uploaded")
		return
	}

	_, err = s.db.Write(hash, r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Error(w, err.Error())
		return
	}

	JSON(w, meta)
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

func (s Server) downloadThumbnail(meta *file.File, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename=thumb_"+meta.Name)
	s.sendfile(meta.Thumbnail, w, r)
}

func (s Server) downloadFile(meta *file.File, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename="+meta.Name)
	s.sendfile(meta.Hash, w, r)
}

func (s Server) getThumbnailDataWithSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	meta, err := s.db.FetchBySlug(slug)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Error(w, err.Error())
		return
	}

	s.downloadThumbnail(meta, w, r)
}

func (s Server) getDataWithSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	meta, err := s.db.FetchBySlug(slug)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Error(w, err.Error())
		return
	}

	s.downloadFile(meta, w, r)
}

func (s Server) getData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchByHash(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Error(w, err.Error())
		return
	}

	s.downloadFile(meta, w, r)
}

func (s Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/file/by-hash/{hash}", s.getData)
	r.Get("/file/by-slug/{slug}", s.getDataWithSlug)
	r.Get("/file/by-slug/{slug}/thumbnail", s.getThumbnailDataWithSlug)
	r.Post("/file/{hash}", s.postData)

	r.Get("/meta/by-hash/{hash}", s.getMeta)
	r.Post("/meta", s.postMeta)

	return r
}

func parseMeta(r io.ReadCloser) (*file.File, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	m := file.File{}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

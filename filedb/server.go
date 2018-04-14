package filedb

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/goware/cors"
	"github.com/zqz/upl/render"
)

// Server implements an HTTP File Uploading Server.
type Server struct {
	db FileDB
}

// NewServer returns a new Server.
func NewServer(db FileDB) Server {
	return Server{
		db: db,
	}
}

func (s Server) postMeta(w http.ResponseWriter, r *http.Request) {
	var m *Meta
	var err error

	if m, err = parseMeta(r.Body); err != nil {
		render.Error(w, "failed to read request")
		return
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
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	render.JSON(w, meta)
}

func (s Server) postData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	_, err := s.db.Write(hash, r.Body)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	render.JSON(w, meta)
}

func (s Server) files(w http.ResponseWriter, r *http.Request) {
	m, err := s.db.List(0)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("err", err.Error())
		render.Error(w, "error loading file list")
		return
	}

	b, _ := json.Marshal(&m)

	w.Write(b)

	pusher, ok := w.(http.Pusher)
	if !ok {
		return
	}

	for _, a := range m {
		pusher.Push("/meta/"+a.Hash+"/thumbnail", nil)
	}
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

	rr, err := s.db.p.Get(hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("err:", err.Error())
		return
	}

	io.Copy(w, rr)

	// err := s.p.Get(hash, w)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	render.Error(w, err.Error())
	// 	return
	// }
}

func (s Server) download(meta *Meta, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Content-Disposition", "inline; filename="+meta.Name)
	s.sendfile(meta.Hash, w, r)
}

func (s Server) getDataWithSlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	meta, err := s.db.FetchMetaWithSlug(slug)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	s.download(meta, w, r)
}

func (s Server) getThumbnail(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	m, err := s.db.m.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	if m.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tns, err := s.db.t.FetchThumbnails([]int{m.ID})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	if len(tns) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t := tns[m.ID]

	s.sendfile(t.Hash, w, r)

}

func (s Server) getData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	s.download(meta, w, r)
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

	r.Get("/files", s.files)

	r.Get("/data/{hash}", s.getData)
	r.Get("/d/{slug}", s.getDataWithSlug)
	r.Post("/data/{hash}", s.postData)

	r.Get("/meta/{hash}/thumbnail", s.getThumbnail)

	r.Post("/meta", s.postMeta)
	r.Get("/meta/{hash}", s.getMeta)

	return r
}

func parseMeta(r io.ReadCloser) (*Meta, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	m := Meta{}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

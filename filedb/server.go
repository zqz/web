package filedb

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"github.com/goware/cors"
	"github.com/zqz/upl/render"
)

type Server struct {
	db FileDB
}

func NewServer(db FileDB) Server {
	return Server{
		db: db,
	}
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

func (s Server) PostMeta(w http.ResponseWriter, r *http.Request) {
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

func (s Server) GetMeta(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}

	render.JSON(w, meta)
}

func (s Server) PostData(w http.ResponseWriter, r *http.Request) {
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

func (s Server) Files(w http.ResponseWriter, r *http.Request) {
	// f := make([]File, 0)

	// dbfiles, err := models.Files(con).All()
	// if err != nil {
	// 	fmt.Println("failed to get all files")
	// }

	// for _, df := range dbfiles {
	// 	file := File{
	// 		Name:  df.Name,
	// 		Size:  df.Size,
	// 		Hash:  df.Hash,
	// 		Token: df.Token,
	// 	}

	// 	f = append(f, file)
	// }

	// b, _ := json.Marshal(&f)

	// w.Write(b)
}

func (s Server) download(meta *Meta, w http.ResponseWriter, r *http.Request) {
	etag := meta.Hash
	w.Header().Set("Content-Type", meta.ContentType)
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Disposition", "inline; filename="+meta.Name)

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			// go lib.TrackDownload(f.DB, file.ID, r, true)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	err := s.db.Read(meta.Hash, w)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Error(w, err.Error())
		return
	}
}

func (s Server) FileDownload(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "token")

	meta, err := s.db.FetchMeta(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		spew.Dump(err)
		render.Error(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	s.download(meta, w, r)
}

func (s Server) GetData(w http.ResponseWriter, r *http.Request) {
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

	r.Get("/files", s.Files)
	// r.Get("/file/{token}/download", s.FileDownload)

	r.Get("/data/{hash}", s.GetData)
	r.Post("/data/{hash}", s.PostData)

	r.Post("/meta", s.PostMeta)
	r.Get("/meta/{hash}", s.GetMeta)

	return r
}

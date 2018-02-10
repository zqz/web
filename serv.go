package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"github.com/goware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/zqz/upl/filedb"
)

var con *sqlx.DB
var db filedb.FileDB

func connect(str string) (*sqlx.DB, error) {
	if len(str) == 0 {
		return nil, errors.New("Empty DB string")
	}

	var err error
	if parsedURL, err := pq.ParseURL(str); err == nil && parsedURL != "" {
		str = parsedURL
	}

	var con *sqlx.DB
	if con, err = sqlx.Connect("postgres", str); err != nil {
		return nil, err
	}

	if err = con.Ping(); err != nil {
		return nil, err
	}

	return con, nil
}

func renderJSON(w http.ResponseWriter, o interface{}) {
	b, err := json.Marshal(o)

	if err != nil {
		renderError(w, http.StatusInternalServerError, "failed to created json")
		return
	}

	spew.Dump(o)

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

type Error struct {
	Message string `json:"message"`
}

func renderError(w http.ResponseWriter, s int, m string) {
	e := Error{
		Message: m,
	}

	fmt.Println("error: ", s, m)

	renderJSON(w, e)
	w.WriteHeader(s)
}

func parseMeta(r io.ReadCloser) (*filedb.Meta, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	defer r.Close()
	m := filedb.Meta{}
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	return &m, err
}

func postMeta(w http.ResponseWriter, r *http.Request) {
	m, err := parseMeta(r.Body)
	if err != nil {
		renderError(w, http.StatusBadRequest, "failed to read request")
		return
	}

	err = db.StoreMeta(m)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "failed to store meta")
		return
	}

	renderJSON(w, m)
}

func getMeta(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	meta, err := db.FetchMeta(hash)

	if err != nil {
		renderError(w, http.StatusNotFound, err.Error())
		return
	}

	renderJSON(w, meta)
}

func postData(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	_, err := db.Write(hash, r.Body)

	if err != nil {
		renderError(w, http.StatusNotFound, err.Error())
		return
	}

	meta, err := db.FetchMeta(hash)

	if err != nil {
		renderError(w, http.StatusNotFound, err.Error())
		return
	}

	renderJSON(w, meta)
}

func files(w http.ResponseWriter, r *http.Request) {
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

func download(meta *filedb.Meta, w http.ResponseWriter, r *http.Request) {
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

	_, err := db.Read(meta.Hash, w)
	if err != nil {
		renderError(w, http.StatusNotFound, err.Error())
		return
	}
}

func fileDownload(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "token")

	meta, err := db.FetchMeta(hash)
	if err != nil {
		renderError(w, http.StatusNotFound, err.Error())
		return
	}

	download(meta, w, r)
}

func main() {
	// os.Mkdir(tmpPath, 0744)
	// os.Mkdir(finalPath, 0744)

	var err error
	con, err = connect("postgres://localhost:5432/zqz2-dev?sslmode=disable")

	if err != nil {
		fmt.Println("error connecting to db", err)
		return
	}

	//mstore := NewDBFileManager(con)
	//db = NewFileManager(mstore)

	db = filedb.NewFileDB(
		filedb.NewMemoryPersistence(),
		filedb.NewMemoryMetaStorage(),
	)
	// var tmpPath string = "/tmp/zqz/"
	// var finalPath string = "/tmp/final/"

	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/files", files)
	r.Get("/file/{token}/download", fileDownload)
	// r.Get("/file/{hash}", fileStatus)
	r.Post("/meta", postMeta)
	r.Post("/data/{hash}", postData)
	r.Get("/meta/{hash}", getMeta)

	http.ListenAndServe(":3001", r)
}

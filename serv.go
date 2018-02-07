package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi"
	"github.com/goware/cors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/zqz/upl/models"
)

var con *sqlx.DB

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

var tmpPath string = "/tmp/zqz/"
var finalPath string = "/tmp/final/"

var uploads map[string]*Upload
var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func init() {
	uploads = make(map[string]*Upload)
}

type Upload struct {
	name          string
	hash          string
	contentType   string
	token         string
	totalSize     int
	bytesReceived int
}

func (u Upload) uploaded() bool {
	return u.bytesReceived == u.totalSize
}

func (u Upload) tmpPath() string {
	return tmpPath + u.hash
}

func (u Upload) finalPath() string {
	return finalPath + u.hash
}

func (u *Upload) write(data io.Reader) error {
	f, err := os.OpenFile(u.tmpPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	// ignore here. likely only ever an EOF error which is expected.
	i, _ := io.Copy(f, data)
	u.bytesReceived += int(i)

	return nil
}

func read(path string, w io.Writer) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}

type File struct {
	Alias       string    `json:"alias"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Hash        string    `json:"hash"`
	Token       string    `json:"token"`
	ContentType string    `json:"type"`
	Size        int       `json:"size"`
	Date        time.Time `json:"date"`
}

type Error struct {
	Message string `json:"message"`
}

func renderError(w http.ResponseWriter, s int, m string) {
	e := Error{
		Message: m,
	}

	fmt.Println("error: ", s, m)

	b, err := json.Marshal(e)
	if err != nil {
		return
	}

	w.WriteHeader(s)
	w.Write(b)
}

type PreparationRequest struct {
	Name        string `json:"name"`
	Size        int    `json:"size"`
	ContentType string `json:"type"`
	Hash        string `json:"hash"`
}

type PreparationResponse struct {
	Token         string `json:"token"`
	BytesReceived int    `json:"bytes_received"`
}

type StatusResponse struct {
	Name          string `json:"name"`
	ContentType   string `json:"type"`
	Hash          string `json:"hash"`
	Token         string `json:"token"`
	BytesReceived int    `json:"bytes_received"`
	Size          int    `json:"size"`
}

func prepareResponseForHash(hash string) *PreparationResponse {
	if f, err := models.Files(con, qm.Where("hash=?", hash)).One(); err == nil {
		return &PreparationResponse{
			Token:         f.Hash,
			BytesReceived: f.Size,
		}
	}

	if u, ok := uploads[hash]; ok == true {
		return &PreparationResponse{
			Token:         u.hash,
			BytesReceived: u.bytesReceived,
		}
	}

	return nil
}

// create an upload, no data.
func dataMeta(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		renderError(w, http.StatusBadRequest, "Failed to read request")
		return
	}

	defer r.Body.Close()

	p := PreparationRequest{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		renderError(w, http.StatusBadRequest, "Failed to read prepare request")
		return
	}

	pr := prepareResponseForHash(p.Hash)

	if pr == nil {
		u := &Upload{
			// token:     token,
			totalSize:   p.Size,
			name:        p.Name,
			hash:        p.Hash,
			contentType: p.ContentType,
		}

		uploads[p.Hash] = u

		pr = &PreparationResponse{
			Token:         p.Hash,
			BytesReceived: 0,
		}
	}

	b, err = json.Marshal(pr)
	if err != nil {
		renderError(
			w,
			http.StatusInternalServerError,
			"Failed to create prepare response",
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)

	spew.Dump(p)
}

func fileStatus(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	if len(hash) == 0 {
		renderError(w, http.StatusBadRequest, "Failed to read token")
		return
	}

	f, err := models.Files(con, qm.Where("Hash=?", hash)).One()
	if err != nil {
		renderError(w, http.StatusNotFound, "no file matched")
		return
	}

	us := StatusResponse{
		Name:          f.Name,
		ContentType:   f.ContentType,
		Size:          f.Size,
		BytesReceived: f.Size,
		Hash:          f.Hash,
		Token:         f.Token,
	}

	b, err := json.Marshal(us)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "failed to created json")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

var uploadedFiles map[string]bool

func fileUploaded(hash string) bool {
	if _, ok := uploadedFiles[hash]; ok {
		return true
	}

	if f, _ := models.Files(con, qm.Where("hash=?", hash)).One(); f != nil {
		uploadedFiles[hash] = true
		return true
	}

	return false
}

// uploadData
func data(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	if len(hash) == 0 {
		renderError(w, http.StatusBadRequest, "no hash specified")
		return
	}

	if fileUploaded(hash) {
		renderError(w, http.StatusBadRequest, "file already uploaded")
		return
	}

	u, ok := uploads[hash]
	if !ok {
		renderError(w, http.StatusBadRequest, "no file with hash")
		return
	}

	if u.uploaded() {
		fmt.Println("already uploaded")
		return
	}

	if err := u.write(r.Body); err != nil {
		fmt.Println("failed to write file", err)
	}

	fmt.Println("received total:", u.bytesReceived)

	if !u.uploaded() {
		return
	}

	fmt.Println("finished")
	os.Rename(u.tmpPath(), u.finalPath())

	f := models.File{
		Name:        u.name,
		Size:        u.bytesReceived,
		Hash:        u.hash,
		Token:       randStr(5),
		ContentType: u.contentType,
	}

	err := f.Insert(con)
	if err != nil {
		fmt.Println("Failed to insert f", err)
		return
	}

	delete(uploads, u.hash)

	b, err := json.Marshal(f)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Failed to build response")
	}

	fmt.Println(string(b))

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func files(w http.ResponseWriter, r *http.Request) {
	f := make([]File, 0)

	dbfiles, err := models.Files(con).All()
	if err != nil {
		fmt.Println("failed to get all files")
	}

	for _, df := range dbfiles {
		file := File{
			Name:  df.Name,
			Size:  df.Size,
			Hash:  df.Hash,
			Token: df.Token,
		}

		f = append(f, file)
	}

	b, _ := json.Marshal(&f)

	w.Write(b)
}

func fileDownload(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if len(token) == 0 {
		renderError(w, http.StatusBadRequest, "failed to read token")
		return
	}

	f, err := models.Files(con, qm.Where("Token=?", token)).One()
	if err != nil {
		renderError(w, http.StatusNotFound, "no file exists")
		return
	}

	etag := f.Hash
	w.Header().Set("Content-Type", f.ContentType)
	w.Header().Set("Etag", etag)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Disposition", "inline; filename="+f.Name)

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			// go lib.TrackDownload(f.DB, file.ID, r, true)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	read(finalPath+f.Hash, w)
}

func main() {
	if err := os.Mkdir(tmpPath, 0744); err == nil {
		fmt.Println("creating tmp folder")
	}
	if err := os.Mkdir(finalPath, 0744); err == nil {
		fmt.Println("creating final folder")
	}

	var err error
	con, err = connect("postgres://localhost:5432/zqz2-dev?sslmode=disable")

	if err != nil {
		fmt.Println("error connecting to db", err)
		return
	}

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
	r.Get("/file/{hash}", fileStatus)
	r.Post("/data/meta", dataMeta)
	r.Post("/data/{hash}", data)

	http.ListenAndServe(":3001", r)
}

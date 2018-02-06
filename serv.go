package main

import (
	"encoding/json"
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
)

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

func (u Upload) Read(w io.Writer) error {
	f, err := os.Open(u.finalPath())
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
	BytesReceived int    `json:"bytes_received"`
	Size          int    `json:"size"`
}

func prepare(w http.ResponseWriter, r *http.Request) {
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

	// token := randstr(10)
	u, ok := uploads[p.Hash]

	if ok != true {
		u = &Upload{
			// token:     token,
			totalSize:   p.Size,
			name:        p.Name,
			hash:        p.Hash,
			contentType: p.ContentType,
		}

		uploads[p.Hash] = u
	}

	pr := PreparationResponse{
		Token:         u.hash,
		BytesReceived: u.bytesReceived,
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

func uploadStatus(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if len(token) == 0 {
		renderError(w, http.StatusBadRequest, "Failed to read token")
		return
	}

	u, ok := uploads[token]
	if !ok {
		renderError(w, http.StatusNotFound, "no file matched")
		return
	}

	us := StatusResponse{
		Name:          u.name,
		ContentType:   u.contentType,
		Size:          u.totalSize,
		BytesReceived: u.bytesReceived,
		Hash:          u.hash,
	}

	b, err := json.Marshal(us)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "failed to created json")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func upload(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if len(token) == 0 {
		renderError(w, http.StatusBadRequest, "Failed to read token")
		return
	}

	u, ok := uploads[token]
	if !ok {
		renderError(w, http.StatusBadRequest, "Invalid Token")
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

	if u.uploaded() {
		fmt.Println("finished")
		os.Rename(u.tmpPath(), u.finalPath())
	}

	f := File{
		Name:        u.name,
		Size:        u.bytesReceived,
		Date:        time.Now(),
		Hash:        u.hash,
		ContentType: u.contentType,
	}

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
	for _, u := range uploads {
		file := File{
			Name: u.name,
			Size: u.totalSize,
			Hash: u.hash,
		}

		f = append(f, file)
	}

	b, _ := json.Marshal(&f)

	w.Write(b)
}

func main() {
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
	r.Get("/upload/{token}", uploadStatus)
	r.Post("/upload/{token}", upload)
	r.Post("/prepare", prepare)

	http.ListenAndServe(":3001", r)
}

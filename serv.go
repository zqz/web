package main

import (
	"errors"
	"net/http"

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

func main() {
	// os.Mkdir(tmpPath, 0744)
	// os.Mkdir(finalPath, 0744)

	s := filedb.Server{
		db: filedb.NewFileDB(
			filedb.NewMemoryPersistence(),
			filedb.NewMemoryMetaStorage(),
		),
	}

	// var err error
	// con, err = connect("postgres://localhost:5432/zqz2-dev?sslmode=disable")

	// if err != nil {
	// 	fmt.Println("error connecting to db", err)
	// 	return
	// }

	//mstore := NewDBFileManager(con)
	//db = NewFileManager(mstore)
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

	r.Get("/files", s.Files)
	r.Get("/file/{token}/download", s.FileDownload)
	// r.Get("/file/{hash}", fileStatus)
	r.Post("/meta", s.PostMeta)
	r.Post("/data/{hash}", s.PostData)
	r.Get("/meta/{hash}", s.GetMeta)

	http.ListenAndServe(":3001", r)
}

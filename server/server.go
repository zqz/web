package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zqz/upl/filedb"
)

type Server struct {
	config   config
	database *sql.DB
	logger   *log.Logger
}

func Init(path string, l *log.Logger) (Server, error) {
	s := Server{}

	cfg, err := parseConfig(path)
	if err != nil {
		return s, err
	}

	l.Println("Parsed Config")

	db, err := cfg.DBConfig.loadDatabase()
	if err != nil {
		return s, err
	}

	l.Println("Connected to DB")

	s.database = db
	s.config = cfg
	s.logger = l

	return s, nil
}

func (s Server) Close() {
	s.database.Close()
}

func (s Server) Run() error {
	db := s.database

	fdb := filedb.NewServer(
		filedb.NewFileDB(
			filedb.NewDiskPersistence(),
			filedb.NewDBMetaStorage(db),
			filedb.NewDBThumbnailStorage(db),
		),
	)
	// fdb.SetLogger(l)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/api", fdb.Router())
	r.Get("/*", serveIndex)
	serveAssets(r)

	s.logger.Println("Listening for web traffic")
	return http.ListenAndServe(":3001", r)
}

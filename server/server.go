package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zqz/upl/filedb"
)

type Server struct {
	config   config
	database *sql.DB
	logger   *log.Logger
}

func (s Server) Log(x ...interface{}) {
	s.logger.Println(x...)
}

func Init(path string) (Server, error) {
	s := Server{}
	s.logger = log.New(os.Stdout, "", log.LstdFlags)

	s.logger.Println("Starting Server")

	cfg, err := parseConfig(path)
	if err != nil {
		return s, err
	}
	s.logger.Println("Parsed Config")

	db, err := cfg.DBConfig.loadDatabase()
	if err != nil {
		return s, err
	}
	s.logger.Println("Connected to DB")

	s.database = db
	s.config = cfg

	return s, nil
}

func (s Server) Close() {
	s.database.Close()
}

func (s Server) runInsecure(r http.Handler) error {
	listenPort := fmt.Sprintf(":%d", s.config.Port)

	s.logger.Println("[server] listening for HTTP traffic on port", listenPort)

	return http.ListenAndServe(listenPort, r)
}

func (s Server) Run() error {
	fdb := filedb.NewServer(
		filedb.NewFileDB(
			filedb.NewDiskPersistence(),
			filedb.NewDBMetaStorage(s.database),
		),
	)

	r := chi.NewRouter()

	logger := middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: s.logger})
	r.Use(logger)

	r.Mount("/api", fdb.Router())

	s.logger.Println("Listening for web traffic")

	return s.run(r)
}

func (s Server) run(r http.Handler) error {
	return s.runInsecure(r)
}

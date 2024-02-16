package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/zqz/upl/filedb"
)

type Server struct {
	config   config
	database *sql.DB
	logger   *zerolog.Logger
}

func (s Server) Log(x ...interface{}) {
	s.logger.Println(x...)
}

func Init(logger *zerolog.Logger, path string) (Server, error) {
	s := Server{}
	s.logger = logger

	s.logger.Info().Msg("initializing")
	cfg, err := parseConfig(path)
	if err != nil {
		return s, err
	}
	s.logger.Info().Msg("config loaded")

	db, err := cfg.DBConfig.loadDatabase()
	if err != nil {
		return s, err
	}
	s.logger.Info().Msg("connected to database")

	s.database = db
	s.config = cfg

	return s, nil
}

func (s Server) Close() {
	s.database.Close()
}

func (s Server) runInsecure(r http.Handler) error {
	listenPort := fmt.Sprintf(":%d", s.config.Port)

	s.logger.Info().Int("port", s.config.Port).Msg("listening for requests")

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
	r.Use(loggerMiddleware(s.logger))
	r.Mount("/api", fdb.Router())

	return s.run(r)
}

func (s Server) run(r http.Handler) error {
	return s.runInsecure(r)
}

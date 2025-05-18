package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
	"github.com/zqz/web/backend/web/pages"
)

type Server struct {
	config   config
	database *sql.DB
	logger   *zerolog.Logger
	env      string
}

func Init(logger *zerolog.Logger, env string) (Server, error) {
	s := Server{
		logger: logger,
	}

	s.logger.Info().Msg("initializing")
	cfg, err := loadConfig(env)
	if err != nil {
		return s, err
	}
	s.logger.Info().Msg("config loaded")

	db, err := openDatabase(cfg.DatabaseURL)
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
	userStorage := userdb.NewDBUserStorage(s.database)
	userDB := userdb.NewDB(userStorage)

	storage, err := filedb.NewDiskPersistence(s.config.FilesPath)
	if err != nil {
		return err
	}

	fdb := filedb.NewFileDB(
		storage,
		filedb.NewDBMetaStorage(s.database),
	)
	fsrv := filedb.NewServer(fdb)

	r := chi.NewRouter()
	if s.config.isDevelopment() {
		s.logger.Info().Msg("running in development mode")

		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "GET", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}).Handler)
	}
	r.Use(loggerMiddleware(s.logger))
	r.Mount("/api", fsrv.Router())
	r.Mount("/", pages.Router(&userDB, &fdb))

	return s.run(r)
}

func (s Server) run(r http.Handler) error {
	return s.runInsecure(r)
}

package web

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
	"github.com/zqz/web/backend/web/api"
	"github.com/zqz/web/backend/web/middleware"
	"github.com/zqz/web/backend/web/pages"
)

type Server struct {
	config   Config
	database *sql.DB
	logger   *zerolog.Logger
	env      string
}

func Init(logger *zerolog.Logger, env string) (Server, error) {
	s := Server{
		logger: logger,
	}

	s.logger.Info().Msg("initializing")
	cfg, err := LoadConfig(env)
	if err != nil {
		return s, err
	}
	s.logger.Info().Msg("config loaded")

	db, err := OpenDatabase(cfg.DatabaseURL)
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
	userStorage := userdb.NewDBStorage(s.database)
	userDB := userdb.NewDB(userStorage)

	storage, err := filedb.NewDiskPersistence(s.config.FilesPath)
	if err != nil {
		return err
	}

	fdb := filedb.NewFileDB(
		storage,
		filedb.NewDBMetaStorage(s.database),
	)

	fdb.AddProcessor(filedb.NewThumbnailProcessor(128))
	apiServer := api.NewServer(fdb, userDB)

	r := chi.NewRouter()
	r.Use(middleware.Auth(&userDB))
	r.Use(middleware.Logging(s.logger))

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

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	r.Mount("/api", apiServer.Router())
	r.Mount("/", pages.Router(&userDB, &fdb))

	return s.run(r)
}

func (s Server) run(r http.Handler) error {
	return s.runInsecure(r)
}

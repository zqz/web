package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/api"
	"github.com/zqz/web/backend/internal/transport/shared/middleware"
	"github.com/zqz/web/backend/internal/transport/web"
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
	s.setupGoth()

	return http.ListenAndServe(listenPort, r)
}

func (s Server) setupGoth() {
	key := "xyz"         // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30 // 30 days
	isProd := false      // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd
	gothic.Store = store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:3001/auth/callback?provider=google",
		),
	)
}

func (s Server) Run() error {
	userStorage := user.NewDBStorage(s.database)
	userDB := user.NewDB(userStorage)

	storage, err := file.NewDiskPersistence(s.config.FilesPath)
	if err != nil {
		return err
	}

	fdb := file.NewFileDB(
		storage,
		file.NewDBMetaStorage(s.database),
	)

	fdb.AddProcessor(file.NewThumbnailProcessor(128))
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

	r.Mount("/", web.DefaultRoutes(&userDB, &fdb))
	r.Mount("/admin", web.AdminRoutes(&userDB, &fdb))

	return s.run(r)
}

func (s Server) run(r http.Handler) error {
	return s.runInsecure(r)
}

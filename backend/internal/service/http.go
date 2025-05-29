package service

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/rs/zerolog"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/api"
	"github.com/zqz/web/backend/internal/transport/shared/middleware"
	"github.com/zqz/web/backend/internal/transport/web"
	"github.com/zqz/web/backend/templates/pages"
)

type Server struct {
	config   Config
	database *sql.DB
	logger   *zerolog.Logger
	env      string
	UserDB   *user.DB
	FileDB   *file.DB

	Router *chi.Mux
}

func NewServer(logger *zerolog.Logger, env string, fdb file.DB, udb *user.DB) (*Server, error) {
	return nil, nil
}

func (s Server) Close() {
	s.database.Close()
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
	store.Options.SameSite = http.SameSiteDefaultMode

	gothic.Store = store

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_KEY"),
			os.Getenv("GOOGLE_SECRET"),
			"http://localhost:3001/auth/callback?provider=google",
		),
	)
}

func NewProdServer(logger *zerolog.Logger, env string) (Server, error) {
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

	userStorage := user.NewDBStorage(s.database)
	userDB := user.NewDB(userStorage)

	storage, err := file.NewDiskPersistence(s.config.FilesPath)
	if err != nil {
		return s, err
	}

	fdb := file.NewDB(
		storage,
		file.NewDBMetaStorage(s.database),
	)

	fdb.AddProcessor(file.NewThumbnailProcessor(128))

	s.FileDB = &fdb
	s.UserDB = &userDB

	s.SetupRoutes()
	s.setupGoth()

	return s, nil
}

func NewTestServer() Server {
	userStorage := user.NewMemoryStorage()
	userDB := user.NewDB(userStorage)

	fileMetaStorage := file.NewMemoryMetaStorage()
	fileFileStorage := file.NewMemoryPersistence()
	fileDB := file.NewDB(
		&fileFileStorage,
		fileMetaStorage,
	)

	logger := zerolog.New(os.Stdout)
	s := Server{
		logger: &logger,
		FileDB: &fileDB,
		UserDB: &userDB,
	}

	s.SetupRoutes()
	s.setupGoth()

	return s
}

func (s *Server) SetupRoutes() {
	r := chi.NewRouter()
	r.Use(middleware.Auth(s.UserDB))
	r.Use(middleware.Logging(s.logger))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		pages.PageError(errors.New("not found")).Render(r.Context(), w)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		pages.PageError(errors.New("not allowed")).Render(r.Context(), w)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	r.Mount("/", web.DefaultRoutes(s.UserDB, s.FileDB))
	r.Mount("/admin", web.AdminRoutes(s.UserDB, s.FileDB))
	r.Mount("/api", api.Routes(s.UserDB, s.FileDB))

	s.Router = r
}

func (s Server) Run() error {
	listenPort := fmt.Sprintf(":%d", s.config.Port)
	s.logger.Info().Int("port", s.config.Port).Msg("listening for requests")

	return http.ListenAndServe(listenPort, s.Router)
}

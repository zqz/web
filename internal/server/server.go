package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/zqz/web/backend/internal/config"
	"github.com/zqz/web/backend/internal/i18n"
	v1 "github.com/zqz/web/backend/internal/handler/api/v1"
	"github.com/zqz/web/backend/internal/handler/auth"
	"github.com/zqz/web/backend/internal/handler/middleware"
	"github.com/zqz/web/backend/internal/handler/web"
	"github.com/zqz/web/backend/internal/repository"
	"github.com/zqz/web/backend/internal/service"
	"github.com/zqz/web/backend/internal/service/processor"
	"github.com/zqz/web/backend/internal/service/storage"
)

// Server holds the HTTP server and dependencies for explicit shutdown.
type Server struct {
	HTTP *http.Server
	pool *pgxpool.Pool
}

// New builds the HTTP handler and server from config and logger.
// Caller must call Shutdown when done to close the database pool.
func New(ctx context.Context, cfg *config.Config, logger *zerolog.Logger) (*Server, error) {
	pool, err := setupDatabase(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	repo := repository.NewRepository(pool)

	stor, err := storage.NewDiskStorage(cfg.FilesPath)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("storage: %w", err)
	}

	fileSvc := service.NewFileService(repo, stor)
	userSvc := service.NewUserService(repo)

	if cfg.EnableThumbnails {
		fileSvc.AddProcessor(processor.NewThumbnailProcessor(cfg.ThumbnailSize))
		logger.Info().Int("size", cfg.ThumbnailSize).Msg("thumbnail processor enabled")
	}

	templates, err := template.New("").Funcs(template.FuncMap{
		"t":        i18n.TFunc(i18n.DefaultLocale),
		"quotejs":  i18n.QuoteJS,
		"i18nJSON": func() template.JS { return i18n.JSON(i18n.DefaultLocale) },
	}).ParseGlob("./templates/*.html")
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("templates: %w", err)
	}

	router := setupRouter(cfg, logger, repo, fileSvc, userSvc, templates)

	srv := &http.Server{
		Addr:         cfg.Address(),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		HTTP: srv,
		pool: pool,
	}, nil
}

// Shutdown gracefully shuts down the HTTP server and closes the database pool.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.HTTP != nil {
		if err := s.HTTP.Shutdown(ctx); err != nil {
			return err
		}
	}
	if s.pool != nil {
		s.pool.Close()
	}
	return nil
}

func setupDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func setupRouter(cfg *config.Config, logger *zerolog.Logger, repo *repository.Repository, fileSvc *service.FileService, userSvc *service.UserService, templates *template.Template) http.Handler {
	r := chi.NewRouter()

	authHandler := auth.NewAuthHandler(userSvc, logger, cfg)

	filesHandler := web.NewFilesHandler(fileSvc, templates)
	pagesHandler := web.NewPagesHandler(templates, userSvc, fileSvc)
	adminHandler := web.NewAdminHandler(repo, templates)

	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(authHandler.AuthMiddleware)
	r.Use(web.PublicUploadsMiddleware(repo))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Mount("/auth", auth.NewRouter(authHandler))

	fileHandler := v1.NewFileHandler(fileSvc)
	userHandler := v1.NewUserHandler(userSvc, fileSvc)
	apiHandler := v1.NewRouter(fileHandler, userHandler, authHandler)
	r.Mount("/api/v1", middleware.RateLimitAPI(repo, logger)(apiHandler))

	fileServer := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", pagesHandler.Upload)
	r.Get("/files", filesHandler.Page)
	r.Get("/files/list", filesHandler.List)
	r.Get("/files/{slug}", pagesHandler.Edit)
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return authHandler.RequireAdmin(next, http.HandlerFunc(pagesHandler.Forbidden))
		})
		r.Get("/admin", adminHandler.Page)
		r.Post("/admin/settings", adminHandler.UpdateSettings)
		r.Get("/users", pagesHandler.Users)
		r.Get("/users/{id}", pagesHandler.UserFiles)
		r.Post("/users/{id}/ban", pagesHandler.UserSetBan)
		r.Post("/users/{id}/unban", pagesHandler.UserSetBan)
		r.Post("/users/{id}/max-file-size", pagesHandler.UserSetMaxFileSize)
		r.Post("/users/{id}/profile", pagesHandler.UserSetProfile)
	})
	r.Get("/api-docs", pagesHandler.APIDocs)
	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return authHandler.RequireAuth(next, http.HandlerFunc(pagesHandler.Unauthorized))
		})
		r.Get("/user", pagesHandler.Profile)
	})

	r.NotFound(pagesHandler.NotFound)

	return r
}

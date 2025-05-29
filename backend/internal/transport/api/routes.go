package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
)

func Routes(users *user.DB, files *file.DB) http.Handler {
	s := NewServer(*files, *users)
	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/file/by-hash/{hash}", s.getData)
	r.Get("/file/by-slug/{slug}", s.getDataWithSlug)
	r.Get("/file/by-slug/{slug}/thumbnail", s.getThumbnailDataWithSlug)
	r.Post("/file/{hash}", s.postData)

	r.Get("/meta/by-hash/{hash}", s.getMeta)
	r.Post("/meta", s.postMeta)

	return r
}

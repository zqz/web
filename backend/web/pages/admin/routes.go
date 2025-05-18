package admin

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
	"github.com/zqz/web/backend/web/helper"
	"github.com/zqz/web/backend/web/middleware"
)

func Router(users *userdb.DB, db *filedb.FileDB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Auth(users))
	r.Use(middleware.AdminOnly)
	r.Use(middleware.Flash)

	r.Get("/users", templ.Handler(PageUsers(users)).ServeHTTP)

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "id")
		u, _ := users.FindById(userId)

		PageUser(u).Render(r.Context(), w)
	})

	r.Get("/files", func(w http.ResponseWriter, r *http.Request) {
		PageFiles(db).Render(r.Context(), w)
	})

	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		PageFile(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/delete", func(w http.ResponseWriter, r *http.Request) {
		helper.AddFlash(w, r, "file was deleted")

		http.Redirect(w, r, "/admin/files", http.StatusTemporaryRedirect)
	})

	return r
}

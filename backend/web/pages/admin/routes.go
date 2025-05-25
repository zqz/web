package admin

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
	"github.com/zqz/web/backend/web/helper"
	"github.com/zqz/web/backend/web/middleware"
)

func Router(users *userdb.DB, db *filedb.FileDB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AdminOnly)
	r.Use(middleware.Flash)

	r.Get("/users", templ.Handler(PageUsers(users)).ServeHTTP)

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userId := chi.URLParam(r, "id")
		u, _ := users.FindById(userId)

		PageUser(u, db).Render(r.Context(), w)
	})

	r.Get("/files", func(w http.ResponseWriter, r *http.Request) {
		PageFiles(db).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/edit", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		PageEditFile(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		userId := strconv.Itoa(f.UserID)
		u, _ := users.FindById(userId)

		PageFile(f, u).Render(r.Context(), w)
	})

	r.Post("/files/{slug}/process", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)
		spew.Dump(f)

		w.Write([]byte("file was re-processed"))
		db.Process(f)
	})

	r.Post("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		f.Comment = r.FormValue("comment")
		f.Name = r.FormValue("name")
		f.Slug = r.FormValue("slug")
		f.Private = len(r.FormValue("private")) > 0

		err := db.UpdateMeta(f)

		if err == nil {
			helper.AddFlash(w, r, "file was edited")
			http.Redirect(w, r, "/admin/files/"+f.Slug, http.StatusFound)
		} else {
			helper.AddFlash(w, r, "failed to save")
			http.Redirect(w, r, "/admin/files/"+f.Slug+"/edit", http.StatusFound)
		}
	})

	r.Delete("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")

		err := db.DeleteMetaWithSlug(slug)
		if err != nil {
			w.Write([]byte("failed to delete " + err.Error()))
			return
		}

		helper.AddFlash(w, r, "file was deleted")
		w.Write([]byte("Deleted"))
	})

	return r
}

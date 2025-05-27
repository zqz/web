package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/shared/helper"
	"github.com/zqz/web/backend/internal/transport/shared/middleware"
	"github.com/zqz/web/backend/templates/admin"
	"github.com/zqz/web/backend/templates/pages"
)

func AdminRoutes(users *user.DB, db *file.FileDB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AdminOnly)
	r.Use(middleware.Flash)

	r.Get("/users", templ.Handler(admin.PageUsers(users)).ServeHTTP)

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		userIdStr := chi.URLParam(r, "id")
		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		u, _ := users.FindById(userId)
		admin.PageUser(u, db).Render(r.Context(), w)
	})

	r.Get("/files", func(w http.ResponseWriter, r *http.Request) {
		admin.PageFiles(db).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/edit", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, err := db.FetchMetaWithSlug(slug)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		admin.PageEditFile(f).Render(r.Context(), w)
	})

	r.Post("/files/{slug}/process", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, err := db.FetchMetaWithSlug(slug)
		if err != nil {
			w.Write([]byte("error: " + err.Error()))
			return
		}

		err = db.Process(f)
		if err != nil {
			w.Write([]byte("error: " + err.Error()))
			return
		}

		w.Write([]byte("success"))
	})

	r.Post("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		f.Comment = r.FormValue("comment")
		f.Name = r.FormValue("name")
		f.Slug = r.FormValue("slug")

		fmt.Println("prv:", r.FormValue("private"))
		f.Private = len(r.FormValue("private")) > 0

		err := db.UpdateMeta(f)
		if err == nil {
			helper.AddFlash(w, r, "file was edited")
			http.Redirect(w, r, "/files/"+f.Slug, http.StatusFound)
		} else {
			helper.AddFlash(w, r, "failed to save "+err.Error())
			http.Redirect(w, r, "/admin/files/"+f.Slug+"/edit", http.StatusTemporaryRedirect)
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

package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/shared/helper"
	"github.com/zqz/web/backend/internal/transport/shared/middleware"
	"github.com/zqz/web/backend/templates/admin"
	"github.com/zqz/web/backend/templates/pages"
)

func findFileBySlug(r *http.Request, files *file.DB) (*file.File, error) {
	slug := chi.URLParam(r, "slug")
	if len(slug) == 0 {
		return nil, errors.New("blank slug")
	}

	f, err := files.FetchMetaWithSlug(slug)
	if f == nil {
		return nil, err
	}

	return f, nil
}

func findFile(w http.ResponseWriter, r *http.Request, files *file.DB) (*file.File, bool) {
	u, err := findFileBySlug(r, files)

	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		pages.PageError(errors.New("file not found")).Render(r.Context(), w)
		return nil, false
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.PageError(err).Render(r.Context(), w)
		return nil, false
	}

	return u, true
}

func findUserById(r *http.Request, users *user.DB) (*user.User, error) {
	userIdStr := chi.URLParam(r, "id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return nil, err
	}

	u, err := users.FindById(userId)
	if u == nil {
		return nil, err
	}

	return u, nil
}

func findUser(w http.ResponseWriter, r *http.Request, users *user.DB) (*user.User, bool) {
	u, err := findUserById(r, users)

	if u == nil {
		w.WriteHeader(http.StatusNotFound)
		pages.PageError(errors.New("user not found")).Render(r.Context(), w)
		return nil, false
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pages.PageError(err).Render(r.Context(), w)
		return nil, false
	}

	return u, true
}

func AdminRoutes(users *user.DB, db *file.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.AdminOnly)
	r.Use(middleware.Flash)

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		admin.PageUsers(users).Render(r.Context(), w)
	})

	r.Get("/users/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		if u, ok := findUser(w, r, users); ok {
			admin.PageUser(u, db).Render(r.Context(), w)
		}
	})

	r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		if u, ok := findUser(w, r, users); ok {
			admin.PageUser(u, db).Render(r.Context(), w)
		}
	})

	r.Get("/files", func(w http.ResponseWriter, r *http.Request) {
		admin.PageFiles(db).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/edit", func(w http.ResponseWriter, r *http.Request) {
		if f, ok := findFile(w, r, db); ok {
			admin.PageEditFile(f).Render(r.Context(), w)
		}
	})

	r.Post("/files/{slug}/process", func(w http.ResponseWriter, r *http.Request) {
		if f, ok := findFile(w, r, db); ok {
			err := db.Process(f)
			if err != nil {
				w.Write([]byte("error: " + err.Error()))
				return
			}

			w.Write([]byte("success"))
		}
	})

	r.Post("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		if f, ok := findFile(w, r, db); ok {
			if comment := r.FormValue("comment"); len(comment) > 0 {
				f.Comment = comment
			}
			if name := r.FormValue("name"); len(name) > 0 {
				f.Name = name
			}
			if slug := r.FormValue("slug"); len(slug) > 0 {
				f.Slug = slug
			}

			f.Private = len(r.FormValue("private")) > 0

			err := db.UpdateMeta(f)
			if err == nil {
				helper.AddFlash(w, r, "file was edited")
				http.Redirect(w, r, "/files/"+f.Slug, http.StatusFound)
			} else {
				helper.AddFlash(w, r, "failed to save "+err.Error())
				http.Redirect(w, r, "/admin/files/"+f.Slug+"/edit", http.StatusTemporaryRedirect)
			}
		}
	})

	r.Delete("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		if f, ok := findFile(w, r, db); ok {
			err := db.DeleteMetaWithSlug(f.Slug)

			if err != nil {
				w.Write([]byte("failed to delete " + err.Error()))
				return
			}

			helper.AddFlash(w, r, "file was deleted")
			w.Write([]byte("Deleted"))
		}
	})

	return r
}

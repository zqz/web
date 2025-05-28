package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"

	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/shared/helper"
	"github.com/zqz/web/backend/templates/pages"
)

func DefaultRoutes(users *user.DB, db *file.FileDB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		u := helper.GetUser(r)
		if u != nil {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		gothic.BeginAuthHandler(w, r)
	})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		existingUser, err := gothic.CompleteUserAuth(w, r)

		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		providerId := existingUser.UserID

		u, err := users.FindByProviderId(providerId)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		if u != nil {
			gothic.StoreInSession("user_id", strconv.Itoa(u.ID), r, w)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		u = &user.User{}
		u.Name = existingUser.Name
		u.Email = existingUser.Email
		u.Provider = "google"
		u.ProviderID = providerId

		err = users.Create(u)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := db.List(0)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		pages.Home(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/preview", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, err := db.FetchMetaWithSlug(slug)

		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		pages.FilePreview(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, err := db.FetchMetaWithSlug(slug)

		if f == nil {
			w.WriteHeader(http.StatusNotFound)
			pages.PageError(errors.New("not found")).Render(r.Context(), w)
		}

		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		u, err := users.FindById(f.UserID)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		pages.PageFile(f, u).Render(r.Context(), w)
	})

	return r
}

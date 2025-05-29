package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/templates/pages"

	_ "github.com/markbates/goth/providers/google"
)

func loginAs(
	w http.ResponseWriter, r *http.Request,
	users *user.DB, au goth.User) {

	providerId := au.UserID

	// check if user exists in DB
	u, err := users.FindByProviderId(providerId)
	if err != nil {
		// db errors only, hopefully...
		pages.PageError(err).Render(r.Context(), w)
		return
	}

	// no user exists, create one.
	if u == nil {
		u = &user.User{}
		u.Name = au.Name
		u.Email = au.Email
		u.Provider = "google" // for now only google is supported
		u.ProviderID = providerId

		err = users.Create(u)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}
	}

	gothic.StoreInSession("user_id", strconv.Itoa(u.ID), r, w)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func DefaultRoutes(users *user.DB, db *file.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		if au, err := gothic.CompleteUserAuth(w, r); err == nil {
			loginAs(w, r, users, au)
		} else {
			gothic.BeginAuthHandler(w, r)
		}
	})

	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		if au, err := gothic.CompleteUserAuth(w, r); err == nil {
			loginAs(w, r, users, au)
			return
		}

		pages.PageError(errors.New("failed to login")).Render(r.Context(), w)
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
		f, err := db.FetchBySlug(slug)

		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		pages.FilePreview(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, err := db.FetchBySlug(slug)

		if f == nil {
			w.WriteHeader(http.StatusNotFound)
			pages.PageError(errors.New("not found")).Render(r.Context(), w)
		}

		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		u, err := users.FindById(f.UserID.Int)
		if err != nil {
			pages.PageError(err).Render(r.Context(), w)
			return
		}

		pages.PageFile(f, u).Render(r.Context(), w)
	})

	return r
}

package web

import (
	"fmt"
	"net/http"
	"strconv"

	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/zqz/web/backend/internal/domain/file"
	"github.com/zqz/web/backend/internal/domain/user"
	"github.com/zqz/web/backend/internal/transport/shared/helper"
	"github.com/zqz/web/backend/templates/pages"
)

func DefaultRoutes(users *user.DB, db *file.FileDB) *chi.Mux {
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
			fmt.Fprintln(w, err)
			return
		}

		providerId := existingUser.UserID

		u, _ := users.FindByProviderId(providerId)

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
			spew.Dump(err)
			w.Write([]byte("failed to create user"))
			return
		}

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		f, _ := db.List(0)

		pages.Home(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/preview", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		pages.FilePreview(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		u, _ := users.FindById(f.UserID)

		pages.PageFile(f, u).Render(r.Context(), w)
	})

	return r
}

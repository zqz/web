package pages

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

	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/userdb"
	"github.com/zqz/web/backend/web/helper"
	"github.com/zqz/web/backend/web/pages/admin"
)

func Router(users *userdb.DB, db *filedb.FileDB) *chi.Mux {
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
		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		providerId := user.UserID

		u, _ := users.FindByProviderId(providerId)

		if u != nil {
			gothic.StoreInSession("user_id", strconv.Itoa(u.ID), r, w)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		u = &userdb.User{}
		u.Name = user.Name
		u.Email = user.Email
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

		//		f1 := f[0]
		Home(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}/preview", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		FilePreview(f).Render(r.Context(), w)
	})

	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)

		userId := strconv.Itoa(f.UserID)
		u, _ := users.FindById(userId)

		PageFile(f, u).Render(r.Context(), w)
	})

	// r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
	// 	slug := chi.URLParam(r, "slug")
	// 	f, _ := db.FetchMetaWithSlug(slug)
	// 	admin.PageAdminFile(f).Render(r.Context(), w)
	// })
	//
	r.Mount("/admin", admin.Router(users, db))

	return r
}

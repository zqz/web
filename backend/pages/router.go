package pages

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"os"

	"github.com/a-h/templ"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/zqz/web/backend/filedb"
	"github.com/zqz/web/backend/models"
	"github.com/zqz/web/backend/userdb"
)

func isAdmin(u *models.User) bool {
	return u != nil && u.Email == "dylan@johnston.ca"
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := getUserFromContext(r.Context())
		if isAdmin(u) {
			next.ServeHTTP(w, r)
			return
		}

		w.Write([]byte("not admin"))
	})
}

func Auth(db *userdb.UserDB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, err := gothic.GetFromSession("user_id", r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			u, err := db.FindUserById(userId)
			if u == nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "user", u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value("user").(*models.User)

	if ok {
		return user
	}

	return nil
}

func getUser(r *http.Request) *models.User {
	return getUserFromContext(r.Context())
}

func Router(udb *userdb.UserDB, db *filedb.FileDB) *chi.Mux {
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

	r.Use(Auth(udb))

	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		u := getUser(r)
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

		u, _ := udb.FindUserByProviderId(providerId)

		if u != nil {
			gothic.StoreInSession("user_id", strconv.Itoa(u.ID), r, w)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		dbu := &models.User{
			Name:       user.Name,
			Email:      user.Email,
			Provider:   "google",
			ProviderID: providerId,
		}

		err = udb.CreateUser(dbu)
		if err != nil {
			spew.Dump(err)
			w.Write([]byte("failed to create user"))
			return
		}

		spew.Dump(dbu)

		data := fmt.Sprintf("%vv", user)
		w.Write([]byte("looks good" + data))
	})

	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		gothic.Logout(w, r)
		w.Write([]byte("logged out"))
	})

	r.Get("/", templ.Handler(Home()).ServeHTTP)
	r.Get("/files/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		f, _ := db.FetchMetaWithSlug(slug)
		File(f).Render(r.Context(), w)
	})

	admin := chi.NewRouter()
	admin.Use(Auth(udb))
	admin.Use(AdminOnly)
	admin.Get("/users", templ.Handler(Users(udb)).ServeHTTP)
	admin.Get("/files", templ.Handler(Files(db)).ServeHTTP)

	r.Mount("/admin", admin)

	return r
}

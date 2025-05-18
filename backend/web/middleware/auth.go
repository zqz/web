package middleware

import (
	"context"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/zqz/web/backend/userdb"
)

func Auth(users *userdb.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId, err := gothic.GetFromSession("user_id", r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			u, err := users.FindById(userId)
			if u == nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), "user", u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

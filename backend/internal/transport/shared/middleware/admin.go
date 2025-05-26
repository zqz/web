package middleware

import (
	"net/http"

	"github.com/zqz/web/backend/internal/transport/shared/helper"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := helper.GetUserFromContext(r.Context())
		if helper.IsAdmin(u) {
			next.ServeHTTP(w, r)
			return
		}

		w.Write([]byte("not admin"))
	})
}

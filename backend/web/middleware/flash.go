package middleware

import (
	"context"
	"net/http"

	"github.com/zqz/web/backend/web/helper"
)

func Flash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flashes := helper.Flashes(w, r)
		ctx := context.WithValue(r.Context(), "flashes", flashes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

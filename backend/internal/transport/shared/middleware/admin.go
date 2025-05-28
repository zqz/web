package middleware

import (
	"errors"
	"net/http"

	"github.com/zqz/web/backend/internal/transport/shared/helper"
	"github.com/zqz/web/backend/templates/pages"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := helper.GetUserFromContext(r.Context())
		if helper.IsAdmin(u) {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusForbidden)
		pages.PageError(errors.New("Not an admin")).Render(r.Context(), w)
	})
}

package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/zqz/web/backend/web/helper"
)

func Flash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		flashes := helper.Flashes(w, r)
		ctx := context.WithValue(r.Context(), "flashes", flashes)

		fmt.Println("flashes in middleware")
		spew.Dump(flashes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

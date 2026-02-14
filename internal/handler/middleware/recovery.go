package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
)

// Recovery creates a panic recovery middleware
func Recovery(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Log the panic
					logger.Error().
						Interface("panic", err).
						Bytes("stack", debug.Stack()).
						Str("method", r.Method).
						Str("path", r.URL.Path).
						Msg("panic recovered")

					// Return 500 error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

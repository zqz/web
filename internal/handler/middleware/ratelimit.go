package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/time/rate"

	"github.com/zqz/web/backend/internal/repository"
)

const (
	settingAPIRateLimitRPS = "api_rate_limit_rps"
	defaultAPIRateLimitRPS = 10
	rateLimitCacheTTL      = time.Second
)

// RateLimitAPI returns a middleware that rate-limits API requests per client IP.
// The limit (requests per second) is read from site setting "api_rate_limit_rps";
// default is 10, 0 means disabled. The value is cached for 1 second to avoid DB on every request.
func RateLimitAPI(repo *repository.Repository, logger *zerolog.Logger) func(next http.Handler) http.Handler {
	var (
		cachedRPS  int = defaultAPIRateLimitRPS
		cachedAt   time.Time
		cacheMu    sync.Mutex
		limiters   sync.Map // map[string]*rate.Limiter keyed by IP
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Refresh RPS from settings (with cache)
			cacheMu.Lock()
			if time.Since(cachedAt) >= rateLimitCacheTTL {
				if val, err := repo.Settings.Get(r.Context(), settingAPIRateLimitRPS); err == nil && val != "" {
					if n, err := strconv.Atoi(val); err == nil && n >= 0 {
						cachedRPS = n
					}
				} else {
					cachedRPS = defaultAPIRateLimitRPS
				}
				cachedAt = time.Now()
			}
			rps := cachedRPS
			cacheMu.Unlock()

			if rps == 0 {
				next.ServeHTTP(w, r)
				return
			}

			ip := r.RemoteAddr
			if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
				// Take first (client) when behind proxy (e.g. nginx)
				if first := strings.TrimSpace(strings.Split(fwd, ",")[0]); first != "" {
					ip = first
				}
			}

			var limiter *rate.Limiter
			if v, ok := limiters.Load(ip); ok {
				limiter = v.(*rate.Limiter)
			} else {
				burst := rps*2 + 1
				if burst < 2 {
					burst = 2
				}
				limiter = rate.NewLimiter(rate.Limit(rps), burst)
				if v, loaded := limiters.LoadOrStore(ip, limiter); loaded {
					limiter = v.(*rate.Limiter)
				}
			}

			if !limiter.Allow() {
				logger.Debug().Str("ip", ip).Int("rps", rps).Msg("API rate limit exceeded")
				w.Header().Set("Retry-After", "1")
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

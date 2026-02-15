package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeoutForNonUpload(t *testing.T) {
	t.Run("non-upload request has 200ms deadline", func(t *testing.T) {
		var gotDeadline time.Time
		var gotOK bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotDeadline, gotOK = r.Context().Deadline()
			w.WriteHeader(http.StatusOK)
		})

		handler := timeoutForNonUpload(next)
		req := httptest.NewRequest(http.MethodGet, pathAPIV1Files, nil)
		rec := httptest.NewRecorder()
		before := time.Now()

		handler.ServeHTTP(rec, req)

		require.True(t, gotOK, "non-upload request should have a deadline")
		assert.WithinDuration(t, before.Add(nonUploadTimeout), gotDeadline, 50*time.Millisecond)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("upload request has no deadline", func(t *testing.T) {
		var gotOK bool
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, gotOK = r.Context().Deadline()
			w.WriteHeader(http.StatusOK)
		})

		handler := timeoutForNonUpload(next)
		req := httptest.NewRequest(http.MethodPost, pathAPIV1Meta+"abc123hash", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.False(t, gotOK, "upload request should not have a deadline from the timeout middleware")
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

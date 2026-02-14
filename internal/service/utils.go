package service

import (
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generateSlug generates a random slug of the specified length
func generateSlug(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// timeFromPgType converts pgtype.Timestamp to time.Time
func timeFromPgType(t pgtype.Timestamp) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

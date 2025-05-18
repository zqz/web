package helper

import (
	"context"
	"net/http"

	"github.com/zqz/web/backend/models"
)

func GetUserFromContext(ctx context.Context) *models.User {
	user, ok := ctx.Value("user").(*models.User)

	if ok {
		return user
	}

	return nil
}

func GetUser(r *http.Request) *models.User {
	return GetUserFromContext(r.Context())
}

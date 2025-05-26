package helper

import (
	"context"
	"net/http"

	"github.com/zqz/web/backend/internal/domain/user"
)

func GetUserFromContext(ctx context.Context) *user.User {
	user, ok := ctx.Value("user").(*user.User)

	if ok {
		return user
	}

	return nil
}

func GetUser(r *http.Request) *user.User {
	return GetUserFromContext(r.Context())
}

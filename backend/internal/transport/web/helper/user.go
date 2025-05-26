package helper

import (
	"context"
	"net/http"

	"github.com/zqz/web/backend/userdb"
)

func GetUserFromContext(ctx context.Context) *userdb.User {
	user, ok := ctx.Value("user").(*userdb.User)

	if ok {
		return user
	}

	return nil
}

func GetUser(r *http.Request) *userdb.User {
	return GetUserFromContext(r.Context())
}

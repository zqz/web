package helper

import "github.com/zqz/web/backend/internal/domain/user"

func IsAdmin(u *user.User) bool {
	return u != nil && (u.Email == "dylan@johnston.ca" || u.Email == "qdylanj@gmail.com")
}

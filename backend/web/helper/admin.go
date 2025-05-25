package helper

import "github.com/zqz/web/backend/userdb"

func IsAdmin(u *userdb.User) bool {
	return false
	// return u != nil && (u.Email == "dylan@johnston.ca" || u.Email == "qdylanj@gmail.com")
}

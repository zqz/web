package helper

import "github.com/zqz/web/backend/models"

func IsAdmin(u *models.User) bool {
	return u != nil && u.Email == "dylan@johnston.ca"
}

package userdb

import (
	"github.com/zqz/web/backend/models"
)

type User struct {
	models.User
}

func validUser(u *User) ([]string, bool) {
	issues := make([]string, 0)

	if len(u.Email) < 4 {
		issues = append(issues, "email invalid")
	}

	if len(u.ProviderID) < 1 {
		issues = append(issues, "a provider id is required")
	}

	return issues, len(issues) == 0
}

package userdb

import (
	"errors"
	"strings"

	"github.com/zqz/web/backend/models"
)

type User struct {
	models.User
}

type DB struct {
	p persister
}

type persister interface {
	Create(*User) error
	FindByProviderId(string) (*User, error)
	FindById(string) (*User, error)
	Update(int, User) (User, bool, error)
	List() ([]*User, error)
}

func NewDB(p persister) DB {
	return DB{
		p: p,
	}
}

func (db DB) FindByProviderId(id string) (*User, error) {
	return db.p.FindByProviderId(id)
}

func (db DB) FindById(id string) (*User, error) {
	return db.p.FindById(id)
}

func (db DB) Create(u *User) error {
	if u == nil {
		return errors.New("a user is required")
	}

	if issues, ok := validUser(u); !ok {
		return errors.New(strings.Join(issues, ", "))
	}

	return db.p.Create(u)
}

func (db DB) List() ([]*User, error) {
	return db.p.List()
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

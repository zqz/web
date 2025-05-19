package userdb

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zqz/web/backend/models"
)

type User struct {
	models.User
}

type DBUserStorage struct {
	db *sql.DB
}

func (s *DBUserStorage) Create(u *User) error {
	if err := u.Insert(context.Background(), s.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func (s *DBUserStorage) FindById(id string) (*User, error) {
	user, err := models.Users(qm.Where("id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return &User{*user}, nil
}

func (s *DBUserStorage) FindByProviderId(id string) (*User, error) {
	user, err := models.Users(qm.Where("provider_id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return &User{*user}, nil
}

func (s *DBUserStorage) List() ([]*User, error) {
	dbUsers, err := models.Users().All(context.Background(), s.db)

	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for _, u := range dbUsers {
		users = append(users, &User{*u})
	}

	return users, err
}

func (s *DBUserStorage) Update(id int, u User) (User, bool, error) {
	return User{}, true, nil
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

func NewDBUserStorage(db *sql.DB) *DBUserStorage {
	return &DBUserStorage{
		db: db,
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

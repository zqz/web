package userdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zqz/web/backend/models"
)

type DBUserStorage struct {
	db    *sql.DB
	users []*User
}

func (s *DBUserStorage) CreateUser(u *User) error {
	dbu := models.User{
		Name:       u.Name,
		Email:      u.Email,
		Provider:   u.Provider,
		ProviderID: u.ProviderID,
	}

	if err := dbu.Insert(context.Background(), s.db, boil.Infer()); err != nil {
		fmt.Println("error", err.Error())
		u.ID = dbu.ID
		return err
	}

	return nil
}

func (s *DBUserStorage) FindUserById(id string) (*User, error) {
	user, err := models.Users(qm.Where("id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
	}, nil
}

func (s *DBUserStorage) FindUserByProviderId(id string) (*User, error) {
	user, err := models.Users(qm.Where("provider_id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
	}, nil
}

func (s *DBUserStorage) List() ([]User, error) {
	dbUsers, err := models.Users().All(context.Background(), s.db)

	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbUsers)+1)
	for _, u := range dbUsers {
		users = append(users, User{
			ID:         u.ID,
			Name:       u.Name,
			Email:      u.Email,
			Provider:   u.Provider,
			ProviderID: u.ProviderID,
		})
	}

	return users, err
}

func (s *DBUserStorage) UpdateUser(id int, u User) (User, bool, error) {
	return User{}, true, nil
}

type User struct {
	ID         int
	Name       string
	Email      string
	Provider   string
	ProviderID string
}

type UserDB struct {
	p persister
}

type persister interface {
	CreateUser(*User) error
	FindUserByProviderId(string) (*User, error)
	FindUserById(string) (*User, error)
	UpdateUser(int, User) (User, bool, error)
	List() ([]User, error)
}

func NewUserDB(p persister) UserDB {
	return UserDB{
		p: p,
	}
}

func NewDBUserStorage(db *sql.DB) *DBUserStorage {
	return &DBUserStorage{
		users: make([]*User, 0),
		db:    db,
	}
}

func (db UserDB) FindUserByProviderId(id string) (*User, error) {
	return db.p.FindUserByProviderId(id)
}

func (db UserDB) FindUserById(id string) (*User, error) {
	return db.p.FindUserById(id)
}

func (db UserDB) CreateUser(u *User) error {
	return db.p.CreateUser(u)
}

func (db UserDB) List() ([]User, error) {
	return db.p.List()
}

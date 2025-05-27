package user

import (
	"context"
	"database/sql"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/zqz/web/backend/internal/models"
)

type DBStorage struct {
	db *sql.DB
}

func (s *DBStorage) Create(u *User) error {
	if err := u.Insert(context.Background(), s.db, boil.Infer()); err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) FindById(id int) (*User, error) {
	user, err := models.Users(qm.Where("id=?", id)).One(context.Background(), s.db)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &User{*user}, nil
}

func (s *DBStorage) FindByProviderId(id string) (*User, error) {
	user, err := models.Users(qm.Where("provider_id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return &User{*user}, nil
}

func (s *DBStorage) List() ([]*User, error) {
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

func (s *DBStorage) Update(id int, u User) (User, bool, error) {
	return User{}, true, nil
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

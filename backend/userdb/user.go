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
	users []*models.User
}

func (s *DBUserStorage) Create(u *models.User) error {
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

func (s *DBUserStorage) FindById(id string) (*models.User, error) {
	user, err := models.Users(qm.Where("id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *DBUserStorage) FindByProviderId(id string) (*models.User, error) {
	user, err := models.Users(qm.Where("provider_id=?", id)).One(context.Background(), s.db)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *DBUserStorage) List() ([]models.User, error) {
	dbUsers, err := models.Users().All(context.Background(), s.db)

	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(dbUsers)+1)
	for _, u := range dbUsers {
		users = append(users, *u)
	}

	return users, err
}

func (s *DBUserStorage) Update(id int, u models.User) (models.User, bool, error) {
	return models.User{}, true, nil
}

type DB struct {
	p persister
}

type persister interface {
	Create(*models.User) error
	FindByProviderId(string) (*models.User, error)
	FindById(string) (*models.User, error)
	Update(int, models.User) (models.User, bool, error)
	List() ([]models.User, error)
}

func NewDB(p persister) DB {
	return DB{
		p: p,
	}
}

func NewDBUserStorage(db *sql.DB) *DBUserStorage {
	return &DBUserStorage{
		users: make([]*models.User, 0),
		db:    db,
	}
}

func (db DB) FindByProviderId(id string) (*models.User, error) {
	return db.p.FindByProviderId(id)
}

func (db DB) FindById(id string) (*models.User, error) {
	return db.p.FindById(id)
}

func (db DB) Create(u *models.User) error {
	return db.p.Create(u)
}

func (db DB) List() ([]models.User, error) {
	return db.p.List()
}

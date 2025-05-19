package userdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func memdb() DB {
	return DB{
		p: NewMemoryStorage(),
	}
}

func TestCreateNilUserFails(t *testing.T) {
	db := memdb()
	err := db.Create(nil)

	assert.NotNil(t, err)
}

func TestCreateEmptyUserFails(t *testing.T) {
	db := memdb()
	err := db.Create(&User{})

	assert.NotNil(t, err)
}

func testuser() User {
	u := User{}
	u.Name = "Test"
	u.Email = "test@site.com"
	u.ProviderID = "123"
	u.Provider = "google"
	return u
}

func TestCreateUserSuccess(t *testing.T) {
	db := memdb()
	u := testuser()
	err := db.Create(&u)
	assert.Nil(t, err)
}

func TestEmptyUserListSize(t *testing.T) {
	db := memdb()
	users, err := db.List()
	assert.Nil(t, err)
	assert.Empty(t, users)
}

func TestUserListWithOneEntrySize(t *testing.T) {
	db := memdb()
	u := testuser()
	db.Create(&u)

	users, err := db.List()
	assert.Nil(t, err)
	assert.NotEmpty(t, users)
}

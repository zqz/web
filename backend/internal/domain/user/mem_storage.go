package user

import "errors"

type MemoryStorage struct {
	userIdx int
	users   []*User
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		userIdx: 0,
		users:   make([]*User, 0),
	}
}

func (s *MemoryStorage) Create(u *User) error {
	u.ID = s.userIdx
	s.userIdx++
	s.users = append(s.users, u)
	return nil
}

func (s *MemoryStorage) FindById(id int) (*User, error) {
	for _, u := range s.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, errors.New("no user found")
}

func (s *MemoryStorage) FindByProviderId(id string) (*User, error) {
	for _, u := range s.users {
		if u.ProviderID == id {
			return u, nil
		}
	}

	return nil, errors.New("no user found")
}

func (s *MemoryStorage) List() ([]*User, error) {
	return s.users, nil
}

func (s *MemoryStorage) Update(id int, u User) (User, bool, error) {
	return User{}, false, nil
}

package userdb

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

func (s *MemoryStorage) FindById(id string) (*User, error) {
	return nil, nil
}

func (s *MemoryStorage) FindByProviderId(id string) (*User, error) {
	return nil, nil
}

func (s *MemoryStorage) List() ([]*User, error) {
	return s.users, nil
}

func (s *MemoryStorage) Update(id int, u User) (User, bool, error) {
	return User{}, false, nil
}

package userstore

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/lattots/piikittaja/pkg/models"
)

type mockStore struct {
	users map[int]*models.User
}

func NewMockStore() (UserStore, error) {
	users := make(map[int]*models.User)
	return &mockStore{users: users}, nil
}

func (s *mockStore) Close() error {
	return nil
}

func (s *mockStore) GetByID(id int) (*models.User, error) {
	user := s.users[id]
	if user == nil {
		return nil, ErrUserNotExists
	}
	return user, nil
}

func (s *mockStore) GetUsers() ([]*models.User, error) {
	return slices.Collect(maps.Values(s.users)), nil
}

func (s *mockStore) SearchUsers(searchTerm string) ([]*models.User, error) {
	result := make([]*models.User, 0)
	for _, u := range s.users {
		if strings.Contains(u.Username, searchTerm) {
			result = append(result, u)
		}
	}
	return result, nil
}

func (s *mockStore) Insert(u *models.User) error {
	exists, err := s.Exists(u.ID)
	if err != nil {
		return fmt.Errorf("error checking if user already exists: %w", err)
	}
	if exists {
		return ErrUserAlreadyExists
	}

	s.users[u.ID] = u
	return nil
}

func (s *mockStore) Update(u *models.User) error {
	exists, err := s.Exists(u.ID)
	if err != nil {
		return fmt.Errorf("error checking if user already exists: %w", err)
	}
	if !exists {
		return ErrUserNotExists
	}

	s.users[u.ID] = u
	return nil
}

func (s *mockStore) Remove(id int) error {
	if _, exists := s.users[id]; exists {
		delete(s.users, id)
	} else {
		return ErrUserNotExists
	}
	return nil
}

func (s *mockStore) Exists(id int) (bool, error) {
	u := s.users[id]
	if u == nil {
		return false, nil
	}
	return true, nil
}

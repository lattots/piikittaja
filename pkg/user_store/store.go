package userstore

import (
	"errors"

	"github.com/lattots/piikittaja/pkg/models"
)

type UserStore interface {
	GetByID(id int) (*models.User, error)

	GetUsers() ([]*models.User, error)
	SearchUsers(searchTerm string) ([]*models.User, error)

	Insert(u *models.User) error
	Update(u *models.User) error
	Remove(id int) error

	Exists(id int) (bool, error)

	Close() error
}

var (
	ErrUserAlreadyExists = errors.New("user already exists in database")
	ErrUserNotExists     = errors.New("user doesn't exist in database")
)

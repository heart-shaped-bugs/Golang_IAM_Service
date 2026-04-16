package repositories

import (
	"errors"
	"iam-service/internal/entities"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	SaveUser(*entities.User) error
	FindByEmail(email string) (*entities.User, error)
}

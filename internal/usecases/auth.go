package auth

import (
	"errors"
	"iam-service/internal/entities"
	"iam-service/internal/repositories"
	hash "iam-service/internal/utils/password"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrEmailIsNotValid    = errors.New("email is not valid")
)

type Service interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
}

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo repositories.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  repo,
		jwtSecret: jwtSecret,
	}
}

func (as *AuthService) Register(email string, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailIsNotValid
	}

	_, err = as.userRepo.FindByEmail(email)
	if err == nil {
		return ErrUserAlreadyExists
	}
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return err
	}

	pwdHash, err := hash.Hash(password)
	if err != nil {
		return err
	}

	userID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	user := &entities.User{
		ID:       userID.String(),
		Email:    email,
		Password: pwdHash,
	}
	err = as.userRepo.SaveUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) Login(email string, password string) (string, error) {
	user, err := as.userRepo.FindByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	err = hash.Compare(user.Password, password)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := as.generateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (as *AuthService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

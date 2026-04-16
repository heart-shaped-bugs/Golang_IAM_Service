package pg_repo

import (
	"database/sql"
	"errors"
	"iam-service/internal/entities"
	"iam-service/internal/repositories"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(user *entities.User) error {
	query := `INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Password,
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	query := `SELECT id, email, password FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) FindUser(ctx context.Context, email, passwordHash string) (*entity.User, error) {
	var user entity.User

	query := `SELECT common_id, roles FROM users WHERE email=$1 AND password_hash=$2`
	err := r.db.GetContext(ctx, &user, query, email, passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("неправильный email или password")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

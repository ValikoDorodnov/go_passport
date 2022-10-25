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

func (r *UserRepository) FindUserByCredentials(ctx context.Context, email, hash string) (*entity.User, error) {
	var user entity.User

	query := `SELECT common_id, roles FROM users WHERE email=$1 AND password_hash=$2`
	err := r.db.GetContext(ctx, &user, query, email, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wrong email or password")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) FindUserById(ctx context.Context, subject string) (*entity.User, error) {
	var user entity.User

	query := `SELECT common_id, roles FROM users WHERE common_id=$1`
	err := r.db.GetContext(ctx, &user, query, subject)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no user")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

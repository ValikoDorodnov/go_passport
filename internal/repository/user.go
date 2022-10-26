package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/jmoiron/sqlx"
)

const (
	findUserByCredentials = `SELECT common_id, roles FROM users WHERE email=$1 AND password_hash=$2`
	findUserBySubject     = `SELECT common_id, roles FROM users WHERE common_id=$1`
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

	err := r.db.GetContext(ctx, &user, findUserByCredentials, email, hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("wrong email or password")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

func (r *UserRepository) FindUserBySubject(ctx context.Context, subject string) (*entity.User, error) {
	var user entity.User

	err := r.db.GetContext(ctx, &user, findUserBySubject, subject)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no user")
		} else {
			return nil, err
		}
	}

	return &user, nil
}

package repository

import (
	"context"
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
	var commonId int
	var roles string

	query := `SELECT common_id, roles FROM users WHERE email=$1 AND password_hash=$2`

	err := r.db.QueryRowContext(ctx, query, email, passwordHash).Scan(&commonId, &roles)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(commonId, roles)
	return user, nil
}

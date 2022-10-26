package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/jmoiron/sqlx"
)

const (
	findByRefresh     = `SELECT subject, platform, expires_in FROM refresh_sessions WHERE refresh_token=$1`
	create            = `INSERT INTO refresh_sessions (subject, refresh_token, platform, expires_in) VALUES ($1, $2, $3, $4)`
	deleteByPlatform  = `DELETE FROM refresh_sessions WHERE subject=$1 AND platform=$2`
	deleteAllSessions = `DELETE FROM refresh_sessions WHERE subject=$1`
)

type RefreshSessionRepository struct {
	db *sqlx.DB
}

func NewRefreshSessionRepository(db *sqlx.DB) *RefreshSessionRepository {
	return &RefreshSessionRepository{
		db: db,
	}
}

func (r *RefreshSessionRepository) FindByRefresh(ctx context.Context, refresh string) (*entity.Session, error) {
	var session entity.Session
	err := r.db.GetContext(ctx, &session, findByRefresh, refresh)
	if err == sql.ErrNoRows {
		return nil, errors.New("no valid session")
	} else {
		return &session, err
	}
}

func (r *RefreshSessionRepository) Create(ctx context.Context, subject string, platform string, token *entity.Token) error {
	res, err := r.db.ExecContext(ctx, create, subject, token.Value, platform, token.Exp)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if int(rows) < 1 {
		return errors.New("failed to create session")
	}

	return err
}

func (r *RefreshSessionRepository) DeleteByPlatform(ctx context.Context, subject string, platform string) error {
	_, err := r.db.ExecContext(ctx, deleteByPlatform, subject, platform)

	return err
}

func (r *RefreshSessionRepository) DeleteAllSessions(ctx context.Context, subject string) error {
	_, err := r.db.ExecContext(ctx, deleteAllSessions, subject)

	return err
}

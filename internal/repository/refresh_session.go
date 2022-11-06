package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/jmoiron/sqlx"
)

const (
	findByRefresh       = `SELECT subject, fingerprint, expires_in FROM refresh_sessions WHERE refresh_token=$1`
	findByFingerprint   = `SELECT 1 FROM refresh_sessions WHERE subject=$1 and fingerprint=$2`
	create              = `INSERT INTO refresh_sessions (subject, refresh_token, fingerprint, expires_in) VALUES ($1, $2, $3, $4)`
	deleteByFingerprint = `DELETE FROM refresh_sessions WHERE subject=$1 AND fingerprint=$2`
	deleteAllSessions   = `DELETE FROM refresh_sessions WHERE subject=$1`
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

func (r *RefreshSessionRepository) FindByFingerprint(ctx context.Context, subject, fingerprint string) error {
	var res bool
	err := r.db.GetContext(ctx, &res, findByFingerprint, subject, fingerprint)
	if err == sql.ErrNoRows {
		return errors.New("no valid session")
	}
	return nil
}

func (r *RefreshSessionRepository) Create(ctx context.Context, subject string, fingerprint string, token *entity.Token) error {
	res, err := r.db.ExecContext(ctx, create, subject, token.Value, fingerprint, token.Exp)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if int(rows) < 1 {
		return errors.New("failed to create session")
	}

	return err
}

func (r *RefreshSessionRepository) DeleteByFingerprint(ctx context.Context, subject string, fingerprint string) error {
	_, err := r.db.ExecContext(ctx, deleteByFingerprint, subject, fingerprint)

	return err
}

func (r *RefreshSessionRepository) DeleteAllSessions(ctx context.Context, subject string) error {
	_, err := r.db.ExecContext(ctx, deleteAllSessions, subject)

	return err
}

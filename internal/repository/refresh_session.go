package repository

import (
	"context"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/go-redis/redis/v9"
)

type RefreshSessionRepository struct {
	redis *redis.Client
}

func NewRefreshSessionRepository(redis *redis.Client) *RefreshSessionRepository {
	return &RefreshSessionRepository{
		redis: redis,
	}
}

func (r *RefreshSessionRepository) CreateSession(ctx context.Context, subject, fingerprint string, token *entity.Token) error {
	err := r.redis.HSet(ctx, token.Value, "subject", subject, "fingerprint", fingerprint).Err()
	if err != nil {
		return err
	}

	err = r.addSessionTtl(ctx, token)
	if err != nil {
		return err
	}

	return r.createSessionPointer(ctx, subject, fingerprint, token)
}

func (r *RefreshSessionRepository) addSessionTtl(ctx context.Context, token *entity.Token) error {
	return r.redis.Expire(ctx, token.Value, token.Exp).Err()
}

func (r *RefreshSessionRepository) createSessionPointer(ctx context.Context, subject, fingerprint string, token *entity.Token) error {
	key := r.key(subject, fingerprint)
	return r.redis.Set(ctx, key, token.Value, token.Exp).Err()
}

func (r *RefreshSessionRepository) FindSession(ctx context.Context, token string) (*entity.Session, error) {
	var session entity.Session
	err := r.redis.HGetAll(ctx, token).Scan(&session)
	if err != nil {
		return nil, err
	} else {
		return &session, err
	}
}

func (r *RefreshSessionRepository) DeleteSessionByPointer(ctx context.Context, subject, fingerprint string) error {
	var err error
	pointer, _ := r.findPointer(ctx, subject, fingerprint)
	if pointer != "" {
		err = r.redis.Del(ctx, pointer).Err()
		if err != nil {
			return err
		}
		err = r.deleteSessionPointer(ctx, subject, fingerprint)
	}

	return err
}

func (r *RefreshSessionRepository) DeleteSessions(ctx context.Context, subject string) error {
	keys, _, err := r.redis.Scan(ctx, 0, subject+"*", 0).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		session, err := r.redis.Get(ctx, key).Result()
		if err == nil {
			err = r.redis.Del(ctx, session).Err()
			err = r.redis.Del(ctx, key).Err()
		}
	}

	return err
}

func (r *RefreshSessionRepository) deleteSessionPointer(ctx context.Context, subject, fingerprint string) error {
	key := r.key(subject, fingerprint)
	return r.redis.Del(ctx, key).Err()
}

func (r *RefreshSessionRepository) findPointer(ctx context.Context, subject, fingerprint string) (string, error) {
	key := r.key(subject, fingerprint)
	return r.redis.Get(ctx, key).Result()
}

func (r *RefreshSessionRepository) key(subject, fingerprint string) string {
	return subject + "_" + fingerprint
}

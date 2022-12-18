package repository

import (
	"context"

	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/go-redis/redis/v9"
)

type SessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(redis *redis.Client) *SessionRepository {
	return &SessionRepository{
		redis: redis,
	}
}

func (r *SessionRepository) Create(ctx context.Context, subject, fingerprint string, token *entity.Token) error {
	err := r.redis.HSet(ctx, token.Value, "subject", subject, "fingerprint", fingerprint).Err()
	if err != nil {
		return err
	}

	err = r.redis.Expire(ctx, token.Value, token.Exp).Err()
	if err != nil {
		return err
	}

	key := r.key(subject, fingerprint)
	return r.redis.Set(ctx, key, token.Value, token.Exp).Err()
}

func (r *SessionRepository) Find(ctx context.Context, token string) (*entity.Session, error) {
	var session entity.Session
	err := r.redis.HGetAll(ctx, token).Scan(&session)
	if err != nil {
		return nil, err
	} else {
		return &session, err
	}
}

func (r *SessionRepository) DeleteByPointer(ctx context.Context, subject, fingerprint string) error {
	var err error
	key := r.key(subject, fingerprint)
	pointer, _ := r.redis.Get(ctx, key).Result()
	if pointer != "" {
		err = r.redis.Del(ctx, pointer).Err()
		if err != nil {
			return err
		}
		err = r.redis.Del(ctx, key).Err()
	}

	return err
}

func (r *SessionRepository) DeleteAll(ctx context.Context, subject string) error {
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

func (r *SessionRepository) CheckTokenIsInBlackList(ctx context.Context, token string) bool {
	_, err := r.redis.Get(ctx, token).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *SessionRepository) AddTokenToBlackList(ctx context.Context, token *entity.ParsedToken) error {
	return r.redis.Set(ctx, token.Jwt, token.Subject, token.ExpTtl).Err()
}

func (r *SessionRepository) key(subject, fingerprint string) string {
	return subject + "_" + fingerprint
}

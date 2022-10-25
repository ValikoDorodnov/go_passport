package repository

import (
	"context"
	"github.com/ValikoDorodnov/go_passport/internal/entity"
	"github.com/go-redis/redis/v9"
)

type AccessSessionRepository struct {
	redis *redis.Client
}

func NewAccessSession(redis *redis.Client) *AccessSessionRepository {
	return &AccessSessionRepository{
		redis: redis,
	}
}

func (r *AccessSessionRepository) CheckTokenIsInBlackList(ctx context.Context, token string) bool {
	_, err := r.redis.Get(ctx, token).Result()
	if err != nil {
		return false
	}
	return true
}

func (r *AccessSessionRepository) AddTokenToBlackList(ctx context.Context, token *entity.ParsedToken) {
	err := r.redis.Set(ctx, token.Jwt, token.Subject, token.ExpTtl).Err()
	if err != nil {
		panic(err)
	}
}

package db

import (
	"fmt"
	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/go-redis/redis/v9"
)

func InitRedis(c config.RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%s", c.Host, c.Port)
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: c.Pass,
		DB:       c.Db,
	})
}

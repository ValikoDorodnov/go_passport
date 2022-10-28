package app

import (
	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http"
	v1 "github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/middleware"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
)

func Configure(conf *config.Config, db *sqlx.DB, redis *redis.Client) *http.Server {
	hash := hasher.NewHasher()

	jwt := service.NewJwtService(conf.Jwt)
	userRepository := repository.NewUserRepository(db)
	sessionRepository := repository.NewRefreshSessionRepository(db)
	accessRepository := repository.NewAccessSession(redis)
	auth := service.NewAuthService(userRepository, sessionRepository, accessRepository, hash, jwt)
	authMiddleware := middleware.NewAuthMiddleware(jwt, accessRepository)

	handler := v1.NewHandler(auth, authMiddleware)
	return http.NewRestServer(conf.Rest, handler.GetRouter())
}

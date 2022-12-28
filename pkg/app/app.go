package app

import (
	"context"
	"fmt"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http"
	v1 "github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/middleware"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/ValikoDorodnov/go_passport/pkg/db"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/ValikoDorodnov/go_passport/pkg/logger"
	v "github.com/ValikoDorodnov/go_passport/pkg/validator"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Server   *http.Server
	postgres *sqlx.DB
	redis    *redis.Client
}

func NewApp(conf *config.Config) (*App, error) {
	p, err := db.InitPostgres(conf.Db)
	if err != nil {
		return nil, err
	}
	r := db.InitRedis(conf.Redis)

	hash := hasher.NewHasher()
	jwt := service.NewJwtService(conf.Jwt)
	userRepository := repository.NewUserRepository(p)
	sessionRepository := repository.NewSessionRepository(r)
	auth := service.NewAuthService(userRepository, sessionRepository, hash, jwt)
	authMiddleware := middleware.NewAuthMiddleware(jwt, sessionRepository)
	validator := v.NewValidator()

	handler := v1.NewHandler(auth, authMiddleware, validator)
	server := http.NewRestServer(conf.Rest, handler.GetRouter())

	return &App{
		postgres: p,
		redis:    r,
		Server:   server,
	}, nil
}

func (r App) Shutdown(ctx context.Context, log *logger.Logger) {
	err := r.redis.Close()
	if err != nil {
		log.Error(fmt.Sprintf("error occured on redis shutting down: %s", err.Error()))
	}
	err = r.postgres.Close()
	if err != nil {
		log.Error(fmt.Sprintf("error occured on postgres shutting down: %s", err.Error()))
	}
	err = r.Server.Shutdown(ctx)
	if err != nil {
		log.Error(fmt.Sprintf("error occured on server shutting down: %s", err.Error()))
	}
}

package app

import (
	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http"
	v1 "github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/middleware"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/ValikoDorodnov/go_passport/pkg/db"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
)

type App struct {
	conf     *config.Config
	postgres *sqlx.DB
	redis    *redis.Client
}

func NewApp(conf *config.Config) (*App, error) {
	p, err := db.InitPostgres(conf.Db)
	if err != nil {
		return nil, err
	}
	r := db.InitRedis(conf.Redis)
	return &App{
		conf:     conf,
		postgres: p,
		redis:    r,
	}, nil
}

func (r App) BuildServer() (*http.Server, error) {
	hash := hasher.NewHasher()
	jwt := service.NewJwtService(r.conf.Jwt)
	userRepository := repository.NewUserRepository(r.postgres)
	sessionRepository := repository.NewRefreshSessionRepository(r.postgres)
	accessRepository := repository.NewAccessSession(r.redis)
	auth := service.NewAuthService(userRepository, sessionRepository, accessRepository, hash, jwt)
	authMiddleware := middleware.NewAuthMiddleware(jwt, accessRepository)

	handler := v1.NewHandler(auth, authMiddleware)
	return http.NewRestServer(r.conf.Rest, handler.GetRouter()), nil
}

func (r App) Shutdown() error {
	err := r.redis.Close()
	if err != nil {
		return err
	}
	return r.postgres.Close()
}

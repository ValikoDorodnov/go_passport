package app

import (
	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http"
	v1 "github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/middleware"
	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/internal/service"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/ValikoDorodnov/go_passport/pkg/logger"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
)

type Srv struct {
	conf  *config.Config
	log   *logger.Logger
	db    *sqlx.DB
	redis *redis.Client
}

func InitSrv(conf *config.Config, log *logger.Logger, db *sqlx.DB, redis *redis.Client) *Srv {
	return &Srv{
		conf:  conf,
		log:   log,
		db:    db,
		redis: redis,
	}
}

func (r Srv) Configure() *http.Server {
	hash := hasher.NewHasher()

	jwt := service.NewJwtService(r.conf.Jwt)
	userRepository := repository.NewUserRepository(r.db)
	sessionRepository := repository.NewRefreshSessionRepository(r.db)
	accessRepository := repository.NewAccessSession(r.redis)
	auth := service.NewAuthService(userRepository, sessionRepository, accessRepository, hash, jwt)
	authMiddleware := middleware.NewAuthMiddleware(jwt, accessRepository)

	handler := v1.NewHandler(auth, authMiddleware)
	return http.NewRestServer(r.conf.Rest, handler.GetRouter())
}

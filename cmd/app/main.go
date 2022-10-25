package main

import (
	"context"
	"fmt"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1/middleware"
	"github.com/ValikoDorodnov/go_passport/pkg/hasher"
	"github.com/ValikoDorodnov/go_passport/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/ValikoDorodnov/go_passport/internal/repository"
	"github.com/ValikoDorodnov/go_passport/pkg/db"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http"
	"github.com/ValikoDorodnov/go_passport/internal/delivery/http/v1"
	"github.com/ValikoDorodnov/go_passport/internal/service"
)

func main() {
	conf := config.InitConfig()
	log := logger.NewLogger()
	hash := hasher.NewHasher()

	postgres, err := db.InitPostgres(conf.Db)
	if err != nil {
		log.Error(fmt.Sprintf("err %v", err))
	}
	defer postgres.Close()

	redis := db.InitRedis(conf.Redis)
	defer redis.Close()

	jwt := service.NewJwtService(conf.Jwt)
	userRepository := repository.NewUserRepository(postgres)
	sessionRepository := repository.NewRefreshSessionRepository(postgres)
	accessRepository := repository.NewAccessSession(redis)
	auth := service.NewAuthService(userRepository, sessionRepository, accessRepository, hash, jwt)
	authMiddleware := middleware.NewAuthMiddleware(jwt, accessRepository)

	handler := v1.NewHandler(auth, authMiddleware)
	srv := http.NewRestServer(conf.Rest, handler.GetRouter())

	go func() {
		fmt.Printf("rest server started at port %s", conf.Rest.Port)
		if err := srv.Run(); err != nil {
			log.Error(fmt.Sprintf("error occured while running http server: %s", err.Error()))
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error(fmt.Sprintf("error occured on server shutting down: %s", err.Error()))
	}
}

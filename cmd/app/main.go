package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ValikoDorodnov/go_passport/pkg/db"
	"github.com/ValikoDorodnov/go_passport/pkg/logger"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/pkg/app"
)

func main() {
	conf := config.InitConfig()
	log := logger.NewLogger()
	postgres, err := db.InitPostgres(conf.Db)
	if err != nil {
		log.Error(fmt.Sprintf("err %v", err))
	}
	defer postgres.Close()

	redis := db.InitRedis(conf.Redis)
	defer redis.Close()

	srv := app.Configure(conf, postgres, redis)

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

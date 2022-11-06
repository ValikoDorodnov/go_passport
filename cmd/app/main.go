package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ValikoDorodnov/go_passport/pkg/logger"

	"github.com/ValikoDorodnov/go_passport/internal/config"
	"github.com/ValikoDorodnov/go_passport/pkg/app"
)

func main() {
	conf := config.InitConfig()
	log := logger.NewLogger()

	application, err := app.NewApp(conf)
	if err != nil {
		log.Error(fmt.Sprintf("error occured while initializating app: %s", err.Error()))
		return
	}

	server, err := application.BuildServer()
	if err != nil {
		log.Error(fmt.Sprintf("error occured while configuring app: %s", err.Error()))
		return
	}

	go func() {
		fmt.Printf("rest server started at port %s", conf.Rest.Port)
		if err := server.Run(); err != nil {
			log.Error(fmt.Sprintf("error occured while running http server: %s", err.Error()))
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := application.Shutdown(); err != nil {
		log.Error(fmt.Sprintf("error occured on databases shutting down: %s", err.Error()))
	}
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(fmt.Sprintf("error occured on server shutting down: %s", err.Error()))
	}
}

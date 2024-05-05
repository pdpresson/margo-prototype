package main

import (
	"gitops_pushservice/appConfig"
	"gitops_pushservice/clients"
	"gitops_pushservice/services"
	"log"
	"os"
	"time"

	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
)

func main() {
	logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	otel.SetLogger(logger)

	config, err := appConfig.InitAppConfig(&logger)
	if err != nil {
		logger.Error(err, "Unable to load app configuration")
		return
	}

	mqClient := clients.NewMQClient(&logger, &config, "state_changed")
	// Give the connection sometime to setup
	<-time.After(time.Second)

	pushService := services.NewPushService(&logger, &config, mqClient)
	pushService.Start()
}

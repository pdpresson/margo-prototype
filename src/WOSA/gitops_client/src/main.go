package main

import (
	"gitops_client/appConfig"
	"gitops_client/clients"
	"gitops_client/services"
	"log"
	"os"

	"github.com/go-logr/stdr"
	"go.opentelemetry.io/otel"
)

func main() {
	logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	otel.SetLogger(logger)

	config, err := appConfig.InitAppConfig(&logger)
	if err != nil {
		logger.Error(err, "Unable to load app configuration")
		panic(err)
	}

	kubectlClient := clients.NewKubectlClient(&logger, &config)
	err = kubectlClient.Init()
	if err != nil {
		logger.Error(err, "Error initializing kubectl client")
		panic(err)
	}
	appService := services.NewAppService(&logger, &config, kubectlClient)
	pullService := services.NewPullService(&logger, &config, appService)
	pullService.Start()
}

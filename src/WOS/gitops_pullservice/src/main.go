package main

import (
	"gitops_pullservice/appConfig"
	"gitops_pullservice/clients"
	"gitops_pullservice/services"
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
		return
	}

	pullConfigClient := clients.NewWosPullConfigClient(&logger, &config)
	appCatalogClient := clients.NewAppCatalogClient(&logger, &config)
	appService := services.NewAppService(&logger, &config, appCatalogClient)
	pullService := services.NewPullService(&logger, &config, pullConfigClient, appService)
	pullService.Start()
}

package main

import (
	"log"
	"orchestration_portal/appConfig"
	"orchestration_portal/server"
	"os"

	"github.com/go-logr/stdr"

	"go.opentelemetry.io/otel"
)

func main() {
	logger := stdr.New(log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile))
	otel.SetLogger(logger)

	logger.Info("Importing configuration")
	config, err := appConfig.InitConfig(&logger, "appConfig")
	if err != nil {
		logger.Error(err, "Unable to load app configuration")
		return
	}

	logger.Info("Initializing HTTP server")
	httpServer := server.InitHttpServer(&logger, &config)
	httpServer.Start()
}

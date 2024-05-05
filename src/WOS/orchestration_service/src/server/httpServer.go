package server

import (
	"fmt"
	"orchestration_service/appConfig"
	"orchestration_service/clients"
	"orchestration_service/controllers"
	"orchestration_service/repositories"
	"orchestration_service/services"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

type HttpServer struct {
	appRepoController *controllers.AppRepoController
	log               *logr.Logger
	router            *gin.Engine
	appConfig         *appConfig.AppConfig
}

func (s HttpServer) Start() {
	s.log.Info("Starting Http Server", s.appConfig.ServerPort)
	err := s.router.Run(s.appConfig.ServerPort)
	if err != nil {
		s.log.Error(err, "Error while starting HTTP server")
	}
}

func InitHttpServer(l *logr.Logger, c *appConfig.AppConfig) HttpServer {
	l.Info("Initializing backend service")
	l.Info(fmt.Sprintf("App config: %+v", c))
	router := gin.Default()

	appRepoRepository := repositories.NewAppRepRepository(l)
	appRepoService := services.NewAppRepoService(l, appRepoRepository)
	appRepoController := controllers.NewAppRepoController(l, appRepoService)

	router.POST(fmt.Sprintf("%s/apprepos", c.RootPath), appRepoController.AddAppRepo)
	router.DELETE(fmt.Sprintf("%s/apprepos/:id", c.RootPath), appRepoController.RemoveAppRepo)
	router.GET(fmt.Sprintf("%s/apprepos", c.RootPath), appRepoController.GetAppRepositories)

	appCatalogRepository := repositories.NewAppCatalogRepository(l)
	appCatalogService := services.NewAppCatalogService(l, appCatalogRepository)
	appCatalogController := controllers.NewAppCatalogController(l, appCatalogService)

	router.GET(fmt.Sprintf("%s/appcatalog", c.RootPath), appCatalogController.GetApps)
	router.GET(fmt.Sprintf("%s/appcatalog/:id/:name/id", c.RootPath), appCatalogController.GetAppId)
	router.GET(fmt.Sprintf("%s/appcatalog/:id", c.RootPath), appCatalogController.GetApp)
	router.POST(fmt.Sprintf("%s/appcatalog/:id/:name", c.RootPath), appCatalogController.AddApp)
	router.PUT(fmt.Sprintf("%s/appcatalog/:id/:name", c.RootPath), appCatalogController.UpdateApp)

	deviceRepoRepository := repositories.NewDeviceRepRepository(l)
	deviceRepoService := services.NewDeviceRepoService(l, c, deviceRepoRepository)
	deviceRepoController := controllers.NewDeviceRepoController(l, deviceRepoService)

	deviceRepository := repositories.NewDeviceRepository(l)
	gitServerClient := clients.NewGetServerClient(l, c)
	mqClient := clients.NewMQClient(l, c, "state_changed")
	// Give the connection sometime to setup
	<-time.After(time.Second)
	deviceService := services.NewDeviceService(l, deviceRepository, deviceRepoService,
		appCatalogService, gitServerClient, mqClient)
	deviceController := controllers.NewDeviceController(l, deviceService, deviceRepoService)

	router.POST(fmt.Sprintf("%s/devices", c.RootPath), deviceController.AddDevice)
	router.GET(fmt.Sprintf("%s/devices", c.RootPath), deviceController.GetDevices)
	router.GET(fmt.Sprintf("%s/devices/:id", c.RootPath), deviceController.GetDevice)
	router.GET(fmt.Sprintf("%s/devices/:id/repo", c.RootPath), deviceRepoController.GetDeviceRepository)
	router.POST(fmt.Sprintf("%s/devices/:id/install/:appId", c.RootPath), deviceController.InstallApp)

	return HttpServer{
		appRepoController: appRepoController,
		log:               l,
		appConfig:         c,
		router:            router,
	}
}

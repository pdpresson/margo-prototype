package server

import (
	"fmt"
	"orchestration_portal/appConfig"
	"orchestration_portal/clients"
	"orchestration_portal/controllers"
	appViews "orchestration_portal/views/appCatalog"
	repoViews "orchestration_portal/views/appRepos"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
)

type HttpServer struct {
	log       *logr.Logger
	router    *gin.Engine
	appConfig *appConfig.AppConfig
}

func (s HttpServer) Start() {
	s.log.Info(fmt.Sprint("Starting Http Server ", s.appConfig.PortalPort))
	err := s.router.Run(s.appConfig.PortalPort)
	if err != nil {
		s.log.Error(err, "Error while starting HTTP server")
	}
}

func InitHttpServer(l *logr.Logger, c *appConfig.AppConfig) HttpServer {
	l.Info("Initializing orchestration portal service")
	l.Info(fmt.Sprintf("App Config: %+v", c))

	router := gin.Default()
	router.LoadHTMLFiles(
		"./views/appRepos/appRepos.gohtml",
		"./views/header.gohtml",
		"./views/appCatalog/appCatalog.gohtml",
		"./views/appCatalog/appInstall.gohtml",
	)

	appReposClient := clients.NewAppReposClient(l, c)
	appReposView := repoViews.NewAppReposView(l)
	appReposController := controllers.NewAppReposController(l, appReposView, appReposClient)

	router.GET(fmt.Sprintf("%s/apprepos", c.RootPath), appReposController.DisplayAppRepos)
	router.POST(fmt.Sprintf("%s/apprepos", c.RootPath), appReposController.AddAppRepo)
	router.DELETE(fmt.Sprintf("%s/apprepos/:id", c.RootPath), appReposController.DeleteAppRepo)

	appCatalogClient := clients.NewAppCatalogClient(l, c)
	deviceClient := clients.NewDeviceClient(l, c)
	appCatalogView := appViews.NewAppCatalogView(l)
	appCatalogController := controllers.NewAppCatalogController(l, appCatalogView, appCatalogClient, deviceClient)
	router.GET(fmt.Sprintf("%s/appcatalog", c.RootPath), appCatalogController.DisplayAppCatalog)
	router.GET(fmt.Sprintf("%s/appcatalog/:id/install", c.RootPath), appCatalogController.DisplayAppInstall)
	router.POST(fmt.Sprintf("%s/appcatalog/:id/install", c.RootPath), appCatalogController.InstallAppOnTarget)

	return HttpServer{
		router:    router,
		log:       l,
		appConfig: c,
	}
}

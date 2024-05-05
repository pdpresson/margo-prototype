package controllers

import (
	"fmt"
	"net/http"
	"orchestration_service/models"
	"orchestration_service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogController struct {
	log               *logr.Logger
	appCatalogService *services.AppCatalogService
}

func (c AppCatalogController) GetApps(ctx *gin.Context) {
	c.log.Info("Adding app")
	response, err := c.appCatalogService.GetApps()
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c AppCatalogController) GetApp(ctx *gin.Context) {
	c.log.Info("Getting app Id")
	c.log.Info(fmt.Sprint("Id param value: ", ctx.Param("id")))
	appId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	response, mErr := c.appCatalogService.GetApp(appId)
	if mErr != nil {
		ctx.JSON(mErr.Status, mErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c AppCatalogController) GetAppId(ctx *gin.Context) {
	c.log.Info("Getting app Id")
	appName := ctx.Param("name")
	repoId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	appId, mErr := c.appCatalogService.GetAppId(repoId, appName)
	if mErr != nil {
		ctx.JSON(mErr.Status, err)
		return
	}

	ctx.JSON(http.StatusOK, appId)
}

func (c AppCatalogController) AddApp(ctx *gin.Context) {
	c.log.Info("Adding app")
	var app models.AppPackage
	if err := ctx.ShouldBindJSON(&app); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, rErr := c.appCatalogService.AddApp(&app)
	if rErr != nil {
		ctx.JSON(rErr.Status, rErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c AppCatalogController) UpdateApp(ctx *gin.Context) {
	c.log.Info("Updating app")
	var app models.AppPackage
	if err := ctx.ShouldBindJSON(&app); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, rErr := c.appCatalogService.UpdateApp(&app)
	if rErr != nil {
		ctx.JSON(rErr.Status, rErr)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func NewAppCatalogController(l *logr.Logger, s *services.AppCatalogService) *AppCatalogController {
	return &AppCatalogController{
		log:               l,
		appCatalogService: s,
	}
}

package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"orchestration_service/models"
	"orchestration_service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppRepoController struct {
	AppRepoService *services.AppRepoService
	log            *logr.Logger
}

func (c AppRepoController) AddAppRepo(ctx *gin.Context) {
	c.log.Info("Adding Repo")
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		c.log.Error(err, "Error while reading add app repository body")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var repo models.AppRepo
	err = json.Unmarshal(body, &repo)
	if err != nil {
		c.log.Error(err, "Error while unmarshalling create result request body")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response, rErr := c.AppRepoService.AddRepository(&repo)
	if rErr != nil {
		ctx.JSON(rErr.Status, rErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c AppRepoController) RemoveAppRepo(ctx *gin.Context) {
	c.log.Info("Remove App Repo")
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		c.log.Error(err, "Error removing the app repo because the id is invalid")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	rErr := c.AppRepoService.DeleteRepository(id)
	if rErr != nil {
		ctx.JSON(rErr.Status, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c AppRepoController) GetAppRepositories(ctx *gin.Context) {
	c.log.Info("Get App Repo")
	response, err := c.AppRepoService.GetAppRepositories()
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func NewAppRepoController(l *logr.Logger, s *services.AppRepoService) *AppRepoController {
	return &AppRepoController{
		log:            l,
		AppRepoService: s,
	}
}

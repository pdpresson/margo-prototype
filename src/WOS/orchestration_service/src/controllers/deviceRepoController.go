package controllers

import (
	"net/http"
	"orchestration_service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type DeviceRepoController struct {
	log               *logr.Logger
	deviceRepoService *services.DeviceRepoService
	deviceService     *services.DeviceService
}

func (c DeviceRepoController) GetDeviceRepository(ctx *gin.Context) {
	c.log.Info("Get Device Repo")
	id, err := uuid.Parse(ctx.Param("Id"))
	if err != nil {
		c.log.Error(err, "Error getting the Device repo because the id is invalid")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	device, mErr := c.deviceService.GetDevice(id)
	if mErr != nil {
		ctx.JSON(mErr.Status, mErr)
		return
	}

	response, mErr := c.deviceRepoService.GetDeviceRepo(device.RepoId)
	if mErr != nil {
		ctx.JSON(mErr.Status, mErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func NewDeviceRepoController(l *logr.Logger, s *services.DeviceRepoService) *DeviceRepoController {
	return &DeviceRepoController{
		log:               l,
		deviceRepoService: s,
	}
}

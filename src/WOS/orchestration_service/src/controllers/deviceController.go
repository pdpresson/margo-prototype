package controllers

import (
	"net/http"
	"orchestration_service/models"
	"orchestration_service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type DeviceController struct {
	log               *logr.Logger
	deviceService     *services.DeviceService
	deviceRepoService *services.DeviceRepoService
}

func (c DeviceController) AddDevice(ctx *gin.Context) {
	c.log.Info("Adding device")
	var device models.DeviceDescription
	if err := ctx.ShouldBindJSON(&device); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	d, rErr := c.deviceService.AddDevice(&device)
	if rErr != nil {
		ctx.JSON(rErr.Status, rErr)
		return
	}

	r, rErr := c.deviceRepoService.GetDeviceRepo(d.RepoId)
	if rErr != nil {
		ctx.JSON(rErr.Status, rErr)
		return
	}

	ctx.JSON(http.StatusOK, struct {
		Device     *models.DeviceDescription
		Repository *models.DeviceRepo
	}{
		Device:     d,
		Repository: r,
	})
}

func (c DeviceController) GetDevice(ctx *gin.Context) {
	c.log.Info("Getting device")
	deviceId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	response, mErr := c.deviceService.GetDevice(deviceId)
	if mErr != nil {
		ctx.JSON(mErr.Status, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c DeviceController) GetDevices(ctx *gin.Context) {
	c.log.Info("Getting devices")
	response, err := c.deviceService.GetDevices()
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (c DeviceController) InstallApp(ctx *gin.Context) {
	c.log.Info("Getting device")
	deviceId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	appId, err := uuid.Parse(ctx.Param("appId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	var properties []models.PropertyValue
	if err := ctx.ShouldBindJSON(&properties); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	mErr := c.deviceService.InstallApp(deviceId, appId, properties)
	if mErr != nil {
		ctx.JSON(mErr.Status, mErr.Error)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

func NewDeviceController(l *logr.Logger, s *services.DeviceService, drs *services.DeviceRepoService) *DeviceController {
	return &DeviceController{
		log:               l,
		deviceService:     s,
		deviceRepoService: drs,
	}
}

package controllers

import (
	"fmt"
	"net/http"
	"orchestration_portal/clients"
	"orchestration_portal/models"
	views "orchestration_portal/views/appCatalog"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogController struct {
	log *logr.Logger

	view          *views.AppCatalogView
	catalogClient *clients.AppCatalogClient
	deviceClient  *clients.DeviceClient
}

func (c AppCatalogController) DisplayAppCatalog(ctx *gin.Context) {
	c.log.Info("Returing app catalog display")
	m := c.catalogClient.GetAppCatalog()
	ctx.HTML(http.StatusOK, c.view.AppCatalogTemplatePath(), c.view.AppCatalogData(m))
}

func (c AppCatalogController) DisplayAppInstall(ctx *gin.Context) {
	c.log.Info("Returning app install display")
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	app := c.catalogClient.GetApp(id)
	devices := c.deviceClient.GetDevices()
	ctx.HTML(http.StatusOK, c.view.AppInstallTemplatePath(), c.view.AppInstallData(app.Description, devices))
}

func (c AppCatalogController) InstallAppOnTarget(ctx *gin.Context) {
	c.log.Info("Installing app on device")

	ctx.Request.ParseForm()

	properties := make([]models.PropertyValue, len(ctx.Request.PostForm)-2)
	idx := -1
	for fieldName, fieldValues := range ctx.Request.PostForm {
		if fieldName == "appid" || fieldName == "devices" {
			continue
		}
		idx++
		for _, value := range fieldValues {
			properties[idx] = models.PropertyValue{
				Name:  fieldName,
				Value: value,
			}
		}
	}

	appId, err := uuid.Parse(ctx.Request.FormValue("appid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	deviceId, err := uuid.Parse(ctx.Request.FormValue("devices"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	c.log.Info(fmt.Sprintf("Installing app %s on device %s", appId, deviceId))
	err = c.deviceClient.AddApp(deviceId, appId, properties)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

}

func NewAppCatalogController(l *logr.Logger, v *views.AppCatalogView,
	ccl *clients.AppCatalogClient, dcl *clients.DeviceClient) *AppCatalogController {

	return &AppCatalogController{
		log:           l,
		view:          v,
		catalogClient: ccl,
		deviceClient:  dcl,
	}
}

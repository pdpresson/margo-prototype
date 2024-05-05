package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"orchestration_portal/appConfig"
	"orchestration_portal/models"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogClient struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

func (c AppCatalogClient) GetAppCatalog() []models.AppDescription {
	c.log.Info("Getting app catalog")

	url, err := url.JoinPath(c.appConfig.ServiceAddress, "appcatalog")
	c.log.Info(fmt.Sprint("Using url ", url))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Unable to create endpoint for host ", c.appConfig.ServiceAddress))
		return nil
	}

	r, err := http.Get(url)
	if err != nil {
		c.log.Error(err, "Unable to get apps")
		return nil
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.log.Error(err, "Unable to read response")
		return nil
	}

	var apps []models.AppDescription
	err = json.Unmarshal(body, &apps)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
	}

	c.log.Info(fmt.Sprint("apps found ", apps))
	return apps
}

func (c AppCatalogClient) GetApp(id uuid.UUID) models.AppPackage {
	c.log.Info("Getting app from catalog")

	route, err := url.JoinPath(c.appConfig.ServiceAddress, "appcatalog", id.String())
	c.log.Info(fmt.Sprint("Using url ", route))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Unable to create endpoint for host ", c.appConfig.ServiceAddress))
	}

	r, err := http.Get(route)
	if err != nil {
		c.log.Error(err, "Unable to get app with id ", id)
		return models.AppPackage{}
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.log.Error(err, "Unable to read response")
		return models.AppPackage{}
	}

	var app models.AppPackage
	err = json.Unmarshal(body, &app)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
		return models.AppPackage{}
	}

	c.log.Info(fmt.Sprint("app found ", app))
	return app
}

func NewAppCatalogClient(l *logr.Logger, c *appConfig.AppConfig) *AppCatalogClient {
	return &AppCatalogClient{
		log:       l,
		appConfig: c,
	}
}

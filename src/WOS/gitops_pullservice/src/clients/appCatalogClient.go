package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitops_pullservice/appConfig"
	"gitops_pullservice/models"
	"io"
	"net/http"
	"net/url"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogClient struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

const (
	routeAppCatalog = "appcatalog"
)

func (c AppCatalogClient) GetAppId(repoId uuid.UUID, appName string) (uuid.UUID, error) {
	c.log.Info(fmt.Sprintf("Getting app ID for app %s in repo ID %s ", appName, repoId.String()))

	url, err := url.JoinPath(c.appConfig.ServiceAddress, routeAppCatalog, repoId.String(), appName, "id")
	if err != nil {
		c.log.Error(err, "Error creating service path")
		return uuid.Nil, err
	}

	r, err := http.Get(url)
	if err != nil {
		c.log.Error(err, "Unable to get repos")
		return uuid.Nil, err
	}
	defer r.Body.Close()

	if r.StatusCode == http.StatusNotFound {
		c.log.Info("App not found")
		return uuid.Nil, err
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.log.Error(err, "Unable to read response")
		return uuid.Nil, err
	}

	id, err := uuid.ParseBytes(body)
	if err != nil {
		c.log.Error(err, "Unable to get id from response")
		return uuid.Nil, nil
	}

	return id, nil
}

func (c AppCatalogClient) UpdateApp(m models.AppPackage) error {
	c.log.Info(fmt.Sprintf("Updating app %s in repo ID %s ", m.Description.Metadata.Name, m.Description.RepoId.String()))

	url, err := url.JoinPath(c.appConfig.ServiceAddress, routeAppCatalog, m.Description.RepoId.String(), m.Description.Metadata.Name)
	if err != nil {
		c.log.Error(err, "Error creating service path")
		return err
	}

	body, err := json.Marshal(m)
	if err != nil {
		c.log.Error(err, "Error marshalling data")
		return err
	}

	r, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Error updating app ", m.Description.Metadata.Name))
		return err
	}

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		c.log.Error(err, fmt.Sprint("Error updating app ", m.Description.Metadata.Name))
	}
	defer res.Body.Close()

	return nil
}

func (c AppCatalogClient) CreateApp(m models.AppPackage) error {
	c.log.Info(fmt.Sprintf("Creating app %s in repo ID %s ", m.Description.Metadata.Name, m.Description.RepoId.String()))

	url, err := url.JoinPath(c.appConfig.ServiceAddress, routeAppCatalog, m.Description.RepoId.String(), m.Description.Metadata.Name)
	if err != nil {
		c.log.Error(err, "Error creating service path")
		return err
	}

	body, err := json.Marshal(m)
	if err != nil {
		c.log.Error(err, "Error marshalling data")
		return err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.log.Error(err, "Unable creating new app")
		return err
	}
	defer res.Body.Close()

	return nil
}

func NewAppCatalogClient(l *logr.Logger, c *appConfig.AppConfig) *AppCatalogClient {
	return &AppCatalogClient{
		log:       l,
		appConfig: c,
	}
}

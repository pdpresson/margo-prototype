package clients

import (
	"bytes"
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

type AppReposClient struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

func (c AppReposClient) GetAppRepos() []models.AppRepo {
	c.log.Info("Getting app repos")

	url, err := url.JoinPath(c.appConfig.ServiceAddress, "apprepos")
	c.log.Info(fmt.Sprint("Using url ", url))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Unable to create endpoint for host ", c.appConfig.ServiceAddress))
	}

	r, err := http.Get(url)
	if err != nil {
		c.log.Error(err, "Unable to get repos")
		return nil
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.log.Error(err, "Unable to read response")
		return nil
	}

	var repos []models.AppRepo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
	}

	c.log.Info(fmt.Sprint("repos found ", repos))
	return repos
}

func (c AppReposClient) DeleteAppRepo(id uuid.UUID) error {
	client := &http.Client{}
	url, err := url.JoinPath(c.appConfig.ServiceAddress, "apprepos", id.String())
	if err != nil {
		c.log.Error(err, "Error creating service url")
		return err
	}

	c.log.Info(fmt.Sprint("Sending delete command to ", url))
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		c.log.Error(err, fmt.Sprint("Error deleting app repo ", id))
		return err
	}
	res, err := client.Do(r)
	if err != nil {
		c.log.Error(err, fmt.Sprint("Error deleting app repo ", id))
	}
	defer res.Body.Close()

	return err
}

func (c AppReposClient) AddAppRep(repoUrl, branch string) error {
	c.log.Info("Getting app repos")

	body, err := json.Marshal(models.AppRepo{
		Url:    repoUrl,
		Branch: branch,
	})
	if err != nil {
		c.log.Error(err, "Error marshalling data")
		return err
	}

	url, err := url.JoinPath(c.appConfig.ServiceAddress, "apprepos")
	c.log.Info(fmt.Sprint("Using url ", url))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Unable to create endpoint for host ", c.appConfig.ServiceAddress))
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.log.Error(err, "Error creating new app repo")
		return err
	}
	defer res.Body.Close()

	return nil
}

func NewAppReposClient(l *logr.Logger, c *appConfig.AppConfig) *AppReposClient {
	return &AppReposClient{
		log:       l,
		appConfig: c,
	}
}

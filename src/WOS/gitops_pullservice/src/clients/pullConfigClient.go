package clients

import (
	"encoding/json"
	"gitops_pullservice/appConfig"
	"gitops_pullservice/models"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-logr/logr"
)

type PullConfigClient struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

func (c PullConfigClient) ReposToMonitor() []*models.AppRepo {
	url, err := url.JoinPath(c.appConfig.ServiceAddress, "apprepos")
	if err != nil {
		c.log.Error(err, "Error creating service path")
		return nil
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

	var repos []*models.AppRepo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
	}

	return repos
}

func (c PullConfigClient) PollFrequency() time.Duration {
	return 60 * time.Second
}

func NewWosPullConfigClient(l *logr.Logger, c *appConfig.AppConfig) *PullConfigClient {
	return &PullConfigClient{
		log:       l,
		appConfig: c,
	}
}

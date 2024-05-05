package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"orchestration_service/appConfig"
	"orchestration_service/models"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type GitServerClient struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

func (c GitServerClient) CreateDeviceRepository(deviceId uuid.UUID) (string, *models.ResponseError) {
	c.log.Info(fmt.Sprint("creating git repository for device ", deviceId))

	existing, mErr := c.GetDeviceReposiotry(deviceId)
	if mErr != nil {
		c.log.Error(errors.New(mErr.Error), "Error trying to check if git repository already exists")
		return "", mErr
	} else if existing != "" {
		c.log.Info("Git repository already for this device already exists")
		return existing, nil
	}

	url, err := url.JoinPath(c.appConfig.DeviceRepoHostName, "/api/v1/admin/users/",
		c.appConfig.DeviceRepoUserName, "/repos")
	if err != nil {
		c.log.Error(err, "error creating repo path")
		return "", &models.ResponseError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
	}
	c.log.Info(fmt.Sprint("Using url ", url))

	payload := fmt.Sprintf(`{"name": "%s", "auto_init": true, "readme":"Default"}`, deviceId)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		c.log.Error(err, "error creating request")
		return "", &models.ResponseError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.appConfig.DeviceRepoToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.log.Error(err, "error sending request")
		return "", &models.ResponseError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated {
		c.log.Error(errors.New("unexpected status code"), fmt.Sprint("received status code ", resp.StatusCode))
		return "", &models.ResponseError{
			Error:  fmt.Sprintf("error creating repository: %s\n", string(body)),
			Status: resp.StatusCode,
		}
	}

	var repo struct {
		CloneUrl string `json:"clone_url"`
	}
	err = json.Unmarshal(body, &repo)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
		return "", &models.ResponseError{
			Error:  fmt.Sprintf("error creating repository: %s\n", string(body)),
			Status: resp.StatusCode,
		}
	}
	repo.CloneUrl = strings.Replace(repo.CloneUrl, ":3030", ":3000", -1)
	return repo.CloneUrl, nil
}

func (c GitServerClient) GetDeviceReposiotry(deviceId uuid.UUID) (string, *models.ResponseError) {
	c.log.Info(fmt.Sprint("Checking if git repository exists for device ", deviceId))
	url, err := url.JoinPath(c.appConfig.DeviceRepoHostName, "/api/v1/repos/",
		c.appConfig.DeviceRepoUserName, deviceId.String())
	if err != nil {
		c.log.Error(err, "error creating repo path")
		return "", &models.ResponseError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
	}
	c.log.Info(fmt.Sprint("Using url ", url))

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		c.log.Error(err, "error creating request")
		return "", &models.ResponseError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", c.appConfig.DeviceRepoToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.log.Error(err, "error sending request")
		return "", &models.ResponseError{
			Error:  err.Error(),
			Status: http.StatusInternalServerError,
		}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	} else if resp.StatusCode == http.StatusOK {
		var repo struct {
			CloneUrl string `json:"clone_url"`
		}
		err = json.Unmarshal(body, &repo)
		if err != nil {
			c.log.Error(err, "Unable to unmarshal response")
			return "", &models.ResponseError{
				Error:  fmt.Sprintf("error creating repository: %s\n", string(body)),
				Status: resp.StatusCode,
			}
		}

		repo.CloneUrl = strings.Replace(repo.CloneUrl, ":3030", ":3000", -1)
		return repo.CloneUrl, nil
	}

	c.log.Error(errors.New("unexpected status code"), fmt.Sprint("received status code ", resp.StatusCode))
	return "", &models.ResponseError{
		Error:  fmt.Sprintf("error creating repository: %s\n", string(body)),
		Status: resp.StatusCode,
	}
}

func (c GitServerClient) getToken() (string, error) {
	c.log.Info("Getting access tokens")
	url, err := url.JoinPath(c.appConfig.DeviceRepoHostName, "/api/v1/users/",
		c.appConfig.DeviceRepoUserName, "/tokens")
	if err != nil {
		c.log.Error(err, "error creating repo path")
		return "", err
	}
	c.log.Info(fmt.Sprint("Using url ", url))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.log.Error(err, "error creating request")
		return "", err
	}
	req.SetBasicAuth(c.appConfig.DeviceRepoUserName, c.appConfig.DeviceRepoPassword)
	resp, err := client.Do(req)
	if err != nil {
		c.log.Error(err, "error sending request")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.log.Error(errors.New("unexpected status code"), fmt.Sprint("received status code ", resp.StatusCode))
		return "", errors.New("unable to get access token")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error(err, "Unable to read response")
		return "", err
	}

	var tokens []struct {
		Name  string `json:"name"`
		Token string `json:"sha1"`
	}
	err = json.Unmarshal(body, &tokens)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
		return "", err
	}

	for _, v := range tokens {
		if v.Name == c.appConfig.DeviceRepoTokenName {
			c.log.Info(fmt.Sprint("Found token ", v.Token))
			return v.Token, nil
		}
	}

	return "", fmt.Errorf("token named %s not found", c.appConfig.DeviceRepoTokenName)
}

func NewGetServerClient(l *logr.Logger, c *appConfig.AppConfig) *GitServerClient {
	return &GitServerClient{
		log:       l,
		appConfig: c,
	}
}

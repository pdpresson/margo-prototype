package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"orchestration_portal/appConfig"
	"orchestration_portal/models"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type DeviceClient struct {
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

func (c DeviceClient) GetDevices() []*models.Device {
	c.log.Info("Getting devices")

	url, err := url.JoinPath(c.appConfig.ServiceAddress, "devices")
	c.log.Info(fmt.Sprint("Using url ", url))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Unable to create endpoint for host ", c.appConfig.ServiceAddress))
		return nil
	}

	r, err := http.Get(url)
	if err != nil {
		c.log.Error(err, "Unable to get devices")
		return nil
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.log.Error(err, "Unable to read response")
		return nil
	}

	var devices []*models.Device
	err = json.Unmarshal(body, &devices)
	if err != nil {
		c.log.Error(err, "Unable to unmarshal response")
		return nil
	}

	c.log.Info(fmt.Sprint("devices found ", devices))
	return devices
}

func (c DeviceClient) AddApp(deviceId, appId uuid.UUID, properties []models.PropertyValue) error {
	c.log.Info("Getting devices")
	url, err := url.JoinPath(c.appConfig.ServiceAddress, "devices", deviceId.String(), "install", appId.String())
	c.log.Info(fmt.Sprint("Using url ", url))
	if err != nil {
		c.log.Error(err, fmt.Sprint("Unable to create endpoint for host ", c.appConfig.ServiceAddress))
		return err
	}

	body, err := json.Marshal(properties)
	if err != nil {
		c.log.Error(err, "unable to marshal properties")
		return err
	}

	r, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.log.Error(err, "Unable to get devices")
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		err := errors.New("unexpected status code")
		c.log.Error(err, fmt.Sprint("Recieved status code ", r.StatusCode))
		return err
	}

	return nil
}

func NewDeviceClient(l *logr.Logger, c *appConfig.AppConfig) *DeviceClient {
	return &DeviceClient{
		log:       l,
		appConfig: c,
	}
}

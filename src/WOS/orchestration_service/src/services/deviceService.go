package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"orchestration_service/clients"
	"orchestration_service/models"
	"orchestration_service/repositories"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type DeviceService struct {
	log               *logr.Logger
	repository        *repositories.DeviceRepository
	deviceRepoService *DeviceRepoService
	appCatalogService *AppCatalogService
	gitClient         *clients.GitServerClient
	mqClient          *clients.MQClient
}

func (s DeviceService) AddDevice(m *models.DeviceDescription) (*models.DeviceDescription, *models.ResponseError) {
	s.log.Info(fmt.Sprintf("Adding device: id='%s' name='%s'", m.ID, m.Metadata.Name))

	if m.ID == uuid.Nil {
		return nil, &models.ResponseError{
			Error:  "Invalid device Id",
			Status: http.StatusBadRequest,
		}
	}

	if m.Metadata.Name == "" {
		return nil, &models.ResponseError{
			Error:  "Invalid device name",
			Status: http.StatusBadRequest,
		}
	}

	gitRepo, err := s.gitClient.CreateDeviceRepository(m.ID)
	if err != nil {
		return nil, err
	}

	result, err := s.deviceRepoService.AddRepository(m.ID, gitRepo)
	if err != nil {
		return nil, err
	}

	m.RepoId = result.ID
	return s.repository.AddDevice(m)
}

func (s DeviceService) GetDevice(id uuid.UUID) (*models.DeviceDescription, *models.ResponseError) {
	s.log.Info(fmt.Sprint("Deleting devices ", id))

	if id == uuid.Nil {
		return nil, &models.ResponseError{
			Error:  "Invalid device Id",
			Status: http.StatusBadRequest,
		}
	}

	return s.repository.GetDevice(id)
}

func (s DeviceService) GetDevices() ([]models.DeviceDescription, *models.ResponseError) {
	s.log.Info("Getting devices")
	return s.repository.GetDevices()
}

func (s DeviceService) InstallApp(deviceId, appId uuid.UUID, properties []models.PropertyValue) *models.ResponseError {
	s.log.Info(fmt.Sprintf("Installing app %s on device %s", appId, deviceId))
	if deviceId == uuid.Nil {
		return &models.ResponseError{
			Error:  "Invalid device Id",
			Status: http.StatusBadRequest,
		}
	}

	if appId == uuid.Nil {
		return &models.ResponseError{
			Error:  "Invalid device Id",
			Status: http.StatusBadRequest,
		}
	}

	app, err := s.appCatalogService.GetApp(appId)
	if err != nil {
		return err
	}

	device, err := s.GetDevice(deviceId)
	if err != nil {
		return err
	}

	deviceRepo, err := s.deviceRepoService.GetDeviceRepo(device.RepoId)
	if err != nil {
		return err
	}

	for _, v := range properties {
		m := app.Description.Properties[v.Name]
		m.Value = v.Value
		app.Description.Properties[v.Name] = m
	}

	appState := models.StateChanged{
		Kind:     models.Install,
		DeviceId: deviceId,
		AppId:    appId,
		AppName:  app.Description.Metadata.Name,
		DeviceRepo: models.DeviceRepo{
			Url:    deviceRepo.Url,
			Branch: deviceRepo.Branch,
		},
		Sources:    app.Description.Sources,
		Properties: app.Description.Properties,
	}

	message, gErr := json.Marshal(appState)
	if gErr != nil {
		return &models.ResponseError{
			Error:  gErr.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	gErr = s.mqClient.Push(message)
	if gErr != nil {
		return &models.ResponseError{
			Error:  gErr.Error(),
			Status: http.StatusInternalServerError,
		}
	}

	err = s.repository.AddApp(deviceId, appId)
	if err != nil {
		return err
	}

	return nil
}

func NewDeviceService(l *logr.Logger, r *repositories.DeviceRepository,
	s *DeviceRepoService, a *AppCatalogService,
	cl *clients.GitServerClient, mcl *clients.MQClient) *DeviceService {
	return &DeviceService{
		log:               l,
		repository:        r,
		deviceRepoService: s,
		appCatalogService: a,
		gitClient:         cl,
		mqClient:          mcl,
	}
}

package services

import (
	"fmt"
	"net/http"
	"orchestration_service/appConfig"
	"orchestration_service/models"
	"orchestration_service/repositories"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type DeviceRepoService struct {
	log        *logr.Logger
	appConfig  *appConfig.AppConfig
	repository *repositories.DeviceRepoRepository
}

func (s DeviceRepoService) AddRepository(deviceId uuid.UUID, url string) (*models.DeviceRepo, *models.ResponseError) {

	if deviceId == uuid.Nil {
		return nil, &models.ResponseError{
			Error:  fmt.Sprintf("'%v' is not a valid device ID", deviceId),
			Status: http.StatusBadRequest,
		}
	}

	repo := models.DeviceRepo{
		Url:    url,
		Branch: "master",
	}

	return s.repository.AddDeviceRepo(&repo)
}

func (s DeviceRepoService) GetDeviceRepo(repoId uuid.UUID) (*models.DeviceRepo, *models.ResponseError) {
	if repoId == uuid.Nil {
		return nil, &models.ResponseError{
			Error:  fmt.Sprintf("'%v' is not a valid ID", repoId),
			Status: http.StatusBadRequest,
		}
	}

	return s.repository.GetDeviceRepo(repoId)
}

func NewDeviceRepoService(l *logr.Logger, c *appConfig.AppConfig, r *repositories.DeviceRepoRepository) *DeviceRepoService {
	return &DeviceRepoService{
		log:        l,
		appConfig:  c,
		repository: r,
	}
}

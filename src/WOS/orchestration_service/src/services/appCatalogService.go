package services

import (
	"fmt"
	"net/http"
	"orchestration_service/models"
	"orchestration_service/repositories"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogService struct {
	log               *logr.Logger
	catalogRepository *repositories.AppCatalogRepository
}

func (s AppCatalogService) GetApps() ([]*models.AppDescription, *models.ResponseError) {
	s.log.Info("Getting apps")
	return s.catalogRepository.GetApps()
}

func (s AppCatalogService) GetApp(id uuid.UUID) (*models.AppPackage, *models.ResponseError) {
	s.log.Info(fmt.Sprint("Getting app with id ", id))

	if id == uuid.Nil {
		return nil, &models.ResponseError{
			Error:  "Invalid app Id",
			Status: http.StatusBadRequest,
		}
	}

	d, err := s.catalogRepository.GetApp(id)
	if err != nil {
		return nil, err
	}

	return &models.AppPackage{
		Description: *d,
	}, nil

}

func (s AppCatalogService) GetAppId(repoId uuid.UUID, appName string) (uuid.UUID, *models.ResponseError) {
	s.log.Info(fmt.Sprintf("Getting app Id for app '%s' in repo '%s'", appName, repoId))

	if repoId == uuid.Nil {
		return uuid.Nil, &models.ResponseError{
			Error:  "Invalid repo Id",
			Status: http.StatusBadRequest,
		}
	}

	if appName == "" {
		return uuid.Nil, &models.ResponseError{
			Error:  "Invalid app name",
			Status: http.StatusBadRequest,
		}
	}

	return s.catalogRepository.GetAppId(repoId, appName)
}

func (s AppCatalogService) AddApp(m *models.AppPackage) (*models.AppPackage, *models.ResponseError) {
	s.log.Info(fmt.Sprintf("Adding app '%s' from repo '%s'", m.Description.Metadata.Name, m.Description.RepoId))

	if m.Description.RepoId == uuid.Nil {
		return m, &models.ResponseError{
			Error:  "Invalid repo ID",
			Status: http.StatusBadRequest,
		}
	}

	if m.Description.Metadata.Name == "" {
		return m, &models.ResponseError{
			Error:  "Invalid app name",
			Status: http.StatusBadRequest,
		}
	}

	d, err := s.catalogRepository.AddApp(&m.Description)
	if err != nil {
		return m, err
	}

	m.Description = *d

	return m, nil
}

func (s AppCatalogService) UpdateApp(m *models.AppPackage) (*models.AppPackage, *models.ResponseError) {
	s.log.Info(fmt.Sprintf("Updating app '%s' with id '%s' from repo '%s'", m.Description.Metadata.Name,
		m.Description.ID, m.Description.RepoId))

	if m.Description.ID == uuid.Nil {
		return m, &models.ResponseError{
			Error:  "Invalid app ID",
			Status: http.StatusBadRequest,
		}
	}

	if m.Description.RepoId == uuid.Nil {
		return m, &models.ResponseError{
			Error:  "Invalid repo ID",
			Status: http.StatusBadRequest,
		}
	}

	if m.Description.Metadata.Name == "" {
		return m, &models.ResponseError{
			Error:  "Invalid app name",
			Status: http.StatusBadRequest,
		}
	}

	d, err := s.catalogRepository.UpdateApp(&m.Description)
	if err != nil {
		return m, err
	}

	m.Description = *d
	return m, nil
}

func NewAppCatalogService(l *logr.Logger, r *repositories.AppCatalogRepository) *AppCatalogService {
	return &AppCatalogService{
		log:               l,
		catalogRepository: r,
	}
}

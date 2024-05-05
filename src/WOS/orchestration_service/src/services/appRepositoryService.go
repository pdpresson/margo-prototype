package services

import (
	"fmt"
	"net/http"
	"net/url"
	"orchestration_service/models"
	"orchestration_service/repositories"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppRepoService struct {
	repository *repositories.AppRepoRepository
	log        *logr.Logger
}

func (s AppRepoService) DeleteRepository(id uuid.UUID) *models.ResponseError {
	if id == uuid.Nil {
		return &models.ResponseError{
			Error:  fmt.Sprintf("'%v' is not a valid ID", id),
			Status: http.StatusBadRequest,
		}
	}

	err := s.repository.DeleteAppRepo(id)
	return err
}

func (s AppRepoService) AddRepository(m *models.AppRepo) (*models.AppRepo, *models.ResponseError) {
	_, err := url.Parse(m.Url)
	if err != nil {
		return nil, &models.ResponseError{
			Error:  fmt.Sprintf("The url '%s' is not valid", m.Url),
			Status: http.StatusBadRequest,
		}
	}

	if m.Branch == "" {
		return nil, &models.ResponseError{
			Error:  fmt.Sprint("A branch must be provided", m.Branch),
			Status: http.StatusBadRequest,
		}
	}

	result, rErr := s.repository.AddAppRepo(m)
	return result, rErr
}

func (s AppRepoService) GetAppRepositories() ([]*models.AppRepo, *models.ResponseError) {
	results, err := s.repository.GetAllAppRepos()
	return results, err
}

func NewAppRepoService(l *logr.Logger, r *repositories.AppRepoRepository) *AppRepoService {
	return &AppRepoService{
		repository: r,
		log:        l,
	}
}

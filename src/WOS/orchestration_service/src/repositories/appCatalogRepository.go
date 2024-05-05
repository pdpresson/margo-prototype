package repositories

import (
	"errors"
	"fmt"
	"net/http"
	"orchestration_service/models"
	"slices"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogRepository struct {
	apps []*models.AppDescription
	log  *logr.Logger
}

func (r *AppCatalogRepository) GetApps() ([]*models.AppDescription, *models.ResponseError) {
	r.log.Info("Getting apps")
	return r.apps, nil
}

func (r *AppCatalogRepository) GetApp(id uuid.UUID) (*models.AppDescription, *models.ResponseError) {
	r.log.Info(fmt.Sprint("Getting app with id ", id))
	idx := slices.IndexFunc(r.apps, func(r *models.AppDescription) bool { return r.ID == id })
	if idx < 0 {
		return nil, &models.ResponseError{
			Error:  fmt.Sprintf("app with id '%s' not found", id),
			Status: http.StatusNotFound,
		}
	}

	return r.apps[idx], nil
}

func (r *AppCatalogRepository) GetAppId(repoId uuid.UUID, appName string) (uuid.UUID, *models.ResponseError) {
	r.log.Info(fmt.Sprintf("Getting id for app '%s' from repo '%s'", appName, repoId))
	idx := slices.IndexFunc(r.apps, func(r *models.AppDescription) bool {
		return r.Metadata.Name == appName && r.RepoId == repoId
	})

	if idx < 0 {
		msg := fmt.Sprintf("Unable to locate app %s from repo %s", appName, repoId)
		r.log.Error(errors.New("application not found"), msg)

		return uuid.Nil, &models.ResponseError{
			Error:  msg,
			Status: http.StatusNotFound,
		}
	}

	return r.apps[idx].ID, nil
}

func (r *AppCatalogRepository) AddApp(m *models.AppDescription) (*models.AppDescription, *models.ResponseError) {
	r.log.Info(fmt.Sprintf("Adding app '%s' from repo '%s'", m.Metadata.Name, m.RepoId))
	idx := slices.IndexFunc(r.apps, func(r *models.AppDescription) bool {
		return r.Metadata.Name == m.Metadata.Name && r.RepoId == m.RepoId
	})

	if idx >= 0 {
		r.log.Info("App already exists")
		return r.apps[idx], nil
	}

	m.ID = uuid.New()
	r.apps = append(r.apps, m)
	return m, nil
}

func (r *AppCatalogRepository) UpdateApp(m *models.AppDescription) (*models.AppDescription, *models.ResponseError) {
	r.log.Info(fmt.Sprint("Updating app ", m.Metadata.Name))
	err := r.deleteApp(m.ID)
	if err != nil {
		return m, err
	}

	r.apps = append(r.apps, m)
	return m, nil
}

func (r *AppCatalogRepository) deleteApp(Id uuid.UUID) *models.ResponseError {
	r.log.Info(fmt.Sprint("Removing app with id ", Id))
	idx := slices.IndexFunc(r.apps, func(r *models.AppDescription) bool { return r.ID == Id })
	if idx < 0 {
		return &models.ResponseError{
			Error:  fmt.Sprintf("Application with id %s not found ", Id),
			Status: http.StatusNotFound,
		}
	}

	r.apps = append(r.apps[:idx], r.apps[idx+1:]...)
	return nil
}

func NewAppCatalogRepository(l *logr.Logger) *AppCatalogRepository {
	return &AppCatalogRepository{
		log: l,
	}
}

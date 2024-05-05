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

type DeviceRepoRepository struct {
	repos []*models.DeviceRepo
	log   *logr.Logger
}

func (r *DeviceRepoRepository) AddDeviceRepo(m *models.DeviceRepo) (*models.DeviceRepo, *models.ResponseError) {
	r.log.Info(fmt.Sprintf("Adding Device repo: url='%s' branch='%s'", m.Url, m.Branch))
	idx := slices.IndexFunc(r.repos, func(r *models.DeviceRepo) bool { return r.Branch == m.Branch && r.Url == m.Url })
	if idx >= 0 {
		r.log.Info("Device repo already exists")
		return r.repos[idx], nil
	}

	m.ID = uuid.New()
	r.repos = append(r.repos, m)
	return m, nil
}

func (r *DeviceRepoRepository) GetDeviceRepo(repoId uuid.UUID) (*models.DeviceRepo, *models.ResponseError) {
	r.log.Info(fmt.Sprint("Getting Device repo ", repoId))
	idx := slices.IndexFunc(r.repos, func(r *models.DeviceRepo) bool { return r.ID == repoId })
	if idx < 0 {
		msg := fmt.Sprint("Unable to locate device repo ", repoId)
		r.log.Error(errors.New("device repo not found"), msg)

		return nil, &models.ResponseError{
			Error:  msg,
			Status: http.StatusNotFound,
		}
	}

	return r.repos[idx], nil
}

func NewDeviceRepRepository(l *logr.Logger) *DeviceRepoRepository {
	return &DeviceRepoRepository{
		log: l,
	}
}

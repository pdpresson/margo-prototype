package repositories

import (
	"fmt"
	"orchestration_service/models"
	"slices"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppRepoRepository struct {
	repos []*models.AppRepo
	log   *logr.Logger
}

func (r *AppRepoRepository) AddAppRepo(m *models.AppRepo) (*models.AppRepo, *models.ResponseError) {
	r.log.Info(fmt.Sprintf("Adding app repo: url='%s' branch='%s'", m.Url, m.Branch))
	idx := slices.IndexFunc(r.repos, func(r *models.AppRepo) bool { return r.Branch == m.Branch && r.Url == m.Url })
	if idx >= 0 {
		r.log.Info("App repo already exists")
		return r.repos[idx], nil
	}

	m.ID = uuid.New()
	r.repos = append(r.repos, m)
	return m, nil
}

func (r *AppRepoRepository) DeleteAppRepo(id uuid.UUID) *models.ResponseError {
	r.log.Info(fmt.Sprint("Removing app repos ", id))
	idx := slices.IndexFunc(r.repos, func(r *models.AppRepo) bool { return r.ID == id })
	r.repos = append(r.repos[:idx], r.repos[idx+1:]...)

	return nil
}

func (r *AppRepoRepository) GetAllAppRepos() ([]*models.AppRepo, *models.ResponseError) {
	r.log.Info("Getting all app repos")
	return r.repos, nil
}

func NewAppRepRepository(l *logr.Logger) *AppRepoRepository {
	return &AppRepoRepository{
		log: l,
	}
}

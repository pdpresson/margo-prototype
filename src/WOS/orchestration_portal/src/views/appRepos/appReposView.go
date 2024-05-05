package views

import (
	"orchestration_portal/models"

	"github.com/go-logr/logr"
)

type AppReposView struct {
	log *logr.Logger
}

type AppReposData struct {
	AppRepos []models.AppRepo
}

func (AppReposView) AppReposTemplatePath() string {
	return "appRepos.gohtml"
}

func (v AppReposView) AppReposData(m []models.AppRepo) *AppReposData {
	return &AppReposData{
		AppRepos: m,
	}
}

func NewAppReposView(l *logr.Logger) *AppReposView {
	return &AppReposView{
		log: l,
	}
}

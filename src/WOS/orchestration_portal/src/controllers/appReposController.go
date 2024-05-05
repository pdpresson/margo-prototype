package controllers

import (
	"fmt"
	"net/http"
	"orchestration_portal/clients"
	views "orchestration_portal/views/appRepos"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppReposController struct {
	client *clients.AppReposClient
	view   *views.AppReposView
	log    *logr.Logger
}

func (c AppReposController) DisplayAppRepos(ctx *gin.Context) {
	c.log.Info("Returing app repos display")
	m := c.client.GetAppRepos()
	ctx.HTML(http.StatusOK, c.view.AppReposTemplatePath(), c.view.AppReposData(m))
}

func (c AppReposController) AddAppRepo(ctx *gin.Context) {
	c.log.Info("Adding new app repo")

	url := ctx.Request.FormValue("url")
	branch := ctx.Request.FormValue("branch")

	err := c.client.AddAppRep(url, branch)
	if err != nil {
		c.log.Error(err, "Error adding new app repository")
	}

	m := c.client.GetAppRepos()
	ctx.HTML(http.StatusCreated, c.view.AppReposTemplatePath(), c.view.AppReposData(m))
}

func (c AppReposController) DeleteAppRepo(ctx *gin.Context) {
	id := ctx.Param("id")
	c.log.Info(fmt.Sprint("Removing App Repo with id ", id))

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.log.Error(err, "Error parsing id from query string")
		m := c.client.GetAppRepos()
		ctx.HTML(http.StatusBadRequest, c.view.AppReposTemplatePath(), c.view.AppReposData(m))
	}

	c.log.Info(fmt.Sprint("Deleting app repo with id ", id))
	err = c.client.DeleteAppRepo(uuid)
	if err != nil {
		c.log.Error(err, fmt.Sprint("Error deleting app repo ", id))
		m := c.client.GetAppRepos()
		ctx.HTML(http.StatusBadRequest, c.view.AppReposTemplatePath(), c.view.AppReposData(m))
	}

	m := c.client.GetAppRepos()
	c.log.Info(fmt.Sprint("App Repos after delete ", m))
	ctx.HTML(http.StatusOK, c.view.AppReposTemplatePath(), c.view.AppReposData(m))
}

func NewAppReposController(l *logr.Logger, v *views.AppReposView,
	cl *clients.AppReposClient) *AppReposController {

	return &AppReposController{
		client: cl,
		view:   v,
		log:    l,
	}
}

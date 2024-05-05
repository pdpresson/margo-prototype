package services

import (
	"fmt"
	"gitops_pullservice/appConfig"
	"gitops_pullservice/clients"
	"gitops_pullservice/models"
	"os"
	"path"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type AppService struct {
	client    *clients.AppCatalogClient
	log       *logr.Logger
	appConfig *appConfig.AppConfig
}

func (s AppService) ProcessRepo(repoId uuid.UUID, repoDir string) error {
	s.log.Info(fmt.Sprint("Processing repo ", repoDir))

	decriptionFile := path.Join(repoDir, s.appConfig.AppDescriptionFileName)
	_, err := os.Stat(decriptionFile)
	if err != nil {
		s.log.Error(err, fmt.Sprintf("Expected file '%s' to exist", decriptionFile))
		return err
	}

	var d models.AppDescription
	data, err := os.ReadFile(decriptionFile)
	if err != nil {
		s.log.Error(err, fmt.Sprint("Error attempting to read ", decriptionFile))
		return err
	}

	err = yaml.Unmarshal(data, &d)
	if err != nil {
		s.log.Error(err, "Error unmarshalling app description")
		return err
	}
	d.RepoId = repoId

	appId, err := s.client.GetAppId(repoId, d.Metadata.Name)
	if err != nil {
		s.log.Error(err, fmt.Sprint("Error getting app id for app ", d.Metadata.Name))
		return err
	}

	appPackage := models.AppPackage{
		Description: d,
	}
	s.log.Info(fmt.Sprint("Created app package ", appPackage))
	if appId == uuid.Nil {
		err = s.client.CreateApp(appPackage)
		if err != nil {
			s.log.Error(err, fmt.Sprint("Error creating app ", d.Metadata.Name))
			return err
		}
		return nil
	}

	d.ID = appId
	err = s.client.UpdateApp(appPackage)
	if err != nil {
		s.log.Error(err, fmt.Sprint("Error updating app", d.Metadata.Name))
	}

	return nil
}

func NewAppService(l *logr.Logger, c *appConfig.AppConfig, cl *clients.AppCatalogClient) *AppService {
	return &AppService{
		client:    cl,
		log:       l,
		appConfig: c,
	}
}

package views

import (
	"orchestration_portal/models"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

type AppCatalogView struct {
	log *logr.Logger
}

type AppCatalogData struct {
	Apps []models.AppDescription
}

type AppInstallData struct {
	ID       uuid.UUID
	Name     string
	Version  string
	Devices  []*models.Device
	Sections []AppInstallPropertySection
}
type AppInstallPropertySection struct {
	Name       string
	Properties []AppInstallProperty
}

type AppInstallProperty struct {
	Property    string
	Name        string
	Description string
}

func (AppCatalogView) AppCatalogTemplatePath() string {
	return "appCatalog.gohtml"
}

func (AppCatalogView) AppCatalogData(m []models.AppDescription) *AppCatalogData {
	return &AppCatalogData{
		Apps: m,
	}
}

func (AppCatalogView) AppInstallTemplatePath() string {
	return "appInstall.gohtml"
}

func (AppCatalogView) AppInstallData(a models.AppDescription, d []*models.Device) *AppInstallData {
	sections := make([]AppInstallPropertySection, len(a.Configuration.Sections))
	for s := range a.Configuration.Sections {
		section := AppInstallPropertySection{
			Name: a.Configuration.Sections[s].Name,
		}

		properties := make([]AppInstallProperty, len(a.Configuration.Sections[s].Settings))
		for i, v := range a.Configuration.Sections[s].Settings {
			properties[i] = AppInstallProperty{
				Property:    v.Property,
				Description: v.Description,
				Name:        v.Name,
			}
		}
		section.Properties = properties
		sections[s] = section
	}

	return &AppInstallData{
		ID:       a.ID,
		Name:     a.Metadata.Name,
		Version:  a.Metadata.Version,
		Devices:  d,
		Sections: sections,
	}
}

func NewAppCatalogView(l *logr.Logger) *AppCatalogView {
	return &AppCatalogView{
		log: l,
	}
}

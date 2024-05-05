package appConfig

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type AppConfig struct {
	RepoRootPath           string `mapstructure:"REPO_ROOT_PATH"`
	AppDescriptionFileName string `mapstructure:"APP_DESCRIPTION_FILE_NAME"`
	ServiceAddress         string `mapstructure:"ORCHESTRATION_SERVICE_ADDRESS"`
}

func InitAppConfig(l *logr.Logger) (AppConfig, error) {
	viper.SetConfigName("appConfig")
	viper.AddConfigPath("./appConfig/")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		return AppConfig{}, err
	}

	return appConfig, nil
}

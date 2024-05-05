package appConfig

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type AppConfig struct {
	RootPath       string `mapstructure:"ORCHESTRATION_PORTAL_ROOT_PATH"`
	PortalPort     string `mapstructure:"ORCHESTRATION_PORTAL_PORT"`
	ServiceAddress string `mapstructure:"ORCHESTRATION_SERVICE_ADDRESS"`
}

func InitConfig(l *logr.Logger, fileName string) (AppConfig, error) {
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

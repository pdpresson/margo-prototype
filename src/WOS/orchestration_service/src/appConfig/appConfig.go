package appConfig

import (
	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type AppConfig struct {
	ServerPort          string `mapstructure:"ORCHESTRATION_SERVICE_PORT"`
	RootPath            string `mapstructure:"ORCHESTRATION_ROOT_PATH"`
	DeviceRepoHostName  string `mapstructure:"DEVICE_REPO_HOSTNAME"`
	DeviceRepoUserName  string `mapstructure:"DEVICE_REPO_USERNAME"`
	DeviceRepoPassword  string `mapstructure:"DEVICE_REPO_PASSWORD"`
	DeviceRepoTokenName string `mapstructure:"DEVICE_REPO_TOKEN_NAME"`
	DeviceRepoToken     string `mapstructure:"DEVICE_REPO_TOKEN"`
	MQAddress           string `mapstructure:"MQ_ADDRESS"`
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

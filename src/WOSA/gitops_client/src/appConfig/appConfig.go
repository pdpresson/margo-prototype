package appConfig

import (
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/spf13/viper"
)

type AppConfig struct {
	DeviceRepoRootPath   string `mapstructure:"DEVICE_REPO_ROOT_PATH"`
	DeviceRepoUrl        string `mapstructure:"DEVICE_REPO_URL"`
	DeviceRepoBranch     string `mapstructure:"DEVICE_REPO_BRANCH"`
	PollFrequency        int    `mapstructure:"POLL_FREQUENCY"`
	DeviceId             string `mapstructure:"DEVICE_ID"`
	CurrentStateRootPath string `mapstructure:"DEVICE_CURRENT_STATE_ROOT_PATH"`
	InCluster            bool   `mapstructure:"IN_CLUSTER"`
}

func (a AppConfig) GetRepoFolder() string {
	return filepath.Join(a.DeviceRepoRootPath, a.DeviceId)
}

func (a AppConfig) GetCurrentStateFolder() string {
	return filepath.Join(a.CurrentStateRootPath, a.DeviceId)
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

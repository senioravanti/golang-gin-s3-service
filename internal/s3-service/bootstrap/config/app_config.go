package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"senioravanti.ru/internal/s3-service/helpers"
)

type ServerConfig struct {
	ReadTimeout time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
}

type AppConfig struct {
	Server ServerConfig `yaml:"server"`
}

func NewAppConfig() (*AppConfig, error) {
	projectRoot, _ := os.Getwd()
	appConfigPath := fmt.Sprintf("%s/configs/s3-service/app.yaml", projectRoot)
	if _, err := os.Stat(appConfigPath); err != nil {
		return nil, helpers.ErrWithCause("config file does not exists", err)
	}

	f, err := os.Open(appConfigPath)
	if err != nil {
		return nil, helpers.ErrWithCause("failed to open config file", err)
	}
	defer f.Close()
	
	var appConfig AppConfig
	if err := yaml.NewDecoder(f).Decode(&appConfig); err != nil {
		return nil, helpers.ErrWithCause("failed to decode config file", err)
	}

	return &appConfig, nil
}
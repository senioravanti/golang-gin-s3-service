package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"senioravanti.ru/internal/s3-service/helpers"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type TlsConfig struct {
	Cert string `yaml:"cert"`
	Key string `yaml:"key"`
}

type ServerConfig struct {
	ReadTimeout time.Duration `yaml:"read-timeout"`
	WriteTimeout time.Duration `yaml:"write-timeout"`
	Port int32 `yaml:"port"`
	Tls *TlsConfig `yaml:"tls"`
}

type S3Config struct {
	AccessKey string `yaml:"access-key"`
	SecretKey string `yaml:"secret-key"`
	Url string `yaml:"url"`
}

type AppConfig struct {
	LogLevel string `yaml:"log-level"`
}

type S3ServiceConfig struct {
	Server *ServerConfig `yaml:"server"`
	S3 *S3Config `yaml:"s3"`
	App *AppConfig `yaml:"app"`
}

func NewS3Client(
	s3Config *S3Config,
) (*s3.Client, error) {
	cfg := aws.Config{}
	cfg.Region = "eu-central-2"
	cfg.Credentials = credentials.NewStaticCredentialsProvider(
		s3Config.AccessKey,
		s3Config.SecretKey,
		"",
	) 
	cfg.BaseEndpoint = aws.String(s3Config.Url)
	
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return s3Client, nil
}

func NewAppConfig() (*S3ServiceConfig, error) {
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
	
	var appConfig S3ServiceConfig
	if err := yaml.NewDecoder(f).Decode(&appConfig); err != nil {
		return nil, helpers.ErrWithCause("failed to decode config file", err)
	}

	return &appConfig, nil
}
package bootstrap

import (
	"fmt"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	"senioravanti.ru/internal/s3-service/helpers"
	"senioravanti.ru/internal/s3-service/bootstrap/config"
)

type Application struct {
	HttpServer *http.Server
	Gin *gin.Engine
	S3Client *s3.Client
}

func loadS3Env() (*config.S3Config, error) {
	s3Env := map[string]string{
		"accessKey": os.Getenv("S3_ACCESS_KEY"),
		"secretKey": os.Getenv("S3_SECRET_KEY"),
		"url": os.Getenv("S3_URL"),
	}

	for key, val := range s3Env {
		if len(val) == 0 {
			return nil, fmt.Errorf("env `%s` is empty", key)
		}
	}

	s3Config := &config.S3Config{
		AccessKey: s3Env["accessKey"],
		SecretKey: s3Env["secretKey"],
		Url: s3Env["url"],
	}
	return s3Config, nil
}

func New() (*Application, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions { AddSource: true, Level: slog.LevelDebug, ReplaceAttr: nil }))
	slog.SetDefault(logger)

	appConfig, err := config.NewAppConfig(); if err != nil {
		slog.Error("failed to load config")
		return nil, err
	}
	
	s3Config, err := loadS3Env()
	if err != nil {
		return nil, fmt.Errorf("failed to load s3 config | %w", err)
	}

	s3Client, err := config.NewS3Client(s3Config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup s3 client | %w", err)
	}

	gin := gin.Default()

	a := &Application {
		HttpServer: &http.Server{
			Addr: "0.0.0.0:" + helpers.Getenv("SERVER_PORT", "8080"),
			ReadTimeout: appConfig.Server.ReadTimeout,
			WriteTimeout: appConfig.Server.WriteTimeout,
			Handler: gin,
		},
		Gin: gin,
		S3Client: s3Client,
	}

	return a, nil 
}

func (a *Application) Run() error {
	return a.HttpServer.ListenAndServe()
}

func (a *Application) Shutdown(ctx context.Context) error {
	return a.HttpServer.Shutdown(ctx)
}
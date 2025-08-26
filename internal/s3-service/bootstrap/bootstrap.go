package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"

	"senioravanti.ru/internal/s3-service/bootstrap/config"
)

type Application struct {
	HttpServer *http.Server
	Tls *config.TlsConfig
	Gin *gin.Engine
	S3Client *s3.Client
}

func newSlog(
	logLevel string,
) {
	var slogLevel slog.Level
	switch logLevel {
		case "ERROR":
			slogLevel = slog.LevelError
		case "WARN":
			slogLevel = slog.LevelWarn
		case "DEBUG":
			slogLevel = slog.LevelDebug
		default:
			slogLevel = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions { 
		AddSource: true, 
		Level: slogLevel, 
		ReplaceAttr: nil, 
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	slog.SetDefault(logger)
}

func New() (*Application, error) {
	s3ServiceConfig, err := config.NewAppConfig() 
	if err != nil {
		slog.Error("failed to load config")
		return nil, err
	}
	newSlog(s3ServiceConfig.App.LogLevel)
	
	s3Client, err := config.NewS3Client(s3ServiceConfig.S3)
	if err != nil {
		return nil, fmt.Errorf("failed to setup s3 client | %w", err)
	}

	gin := gin.Default()

	a := &Application {
		HttpServer: &http.Server{
			Addr: fmt.Sprintf("0.0.0.0:%d", s3ServiceConfig.Server.Port),
			ReadTimeout: s3ServiceConfig.Server.ReadTimeout,
			WriteTimeout: s3ServiceConfig.Server.WriteTimeout,
			Handler: gin,
		},
		Tls: s3ServiceConfig.Server.Tls,
		Gin: gin,
		S3Client: s3Client,
	}

	return a, nil 
}

func (a *Application) Run() error {
	if a.Tls != nil {
		return a.HttpServer.ListenAndServeTLS(a.Tls.Cert, a.Tls.Key)
	}
	return a.HttpServer.ListenAndServe()
}

func (a *Application) Shutdown(ctx context.Context) error {
	return a.HttpServer.Shutdown(ctx)
}
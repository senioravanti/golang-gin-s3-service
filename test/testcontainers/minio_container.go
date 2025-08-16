package testcontainers

import (
	"context"

	"log/slog"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/minio"
)

type MinioContainer struct {
	Container *minio.MinioContainer
	Url string
}

func NewMinioContainer() (*MinioContainer, error) {
	ctx := context.Background()

	container, err := minio.Run(ctx, "minio/minio:latest")
	if err != nil {
		slog.Error("failed to run container")
		return nil, err
	}
	
	minioUrl, err := container.ConnectionString(ctx)
	if err != nil {
		slog.Error("failed to get url")
		return nil, err
	}
	
	minioContainer := &MinioContainer{
		Container: container,
		Url: minioUrl,
	}
	return minioContainer, nil
}

func (m *MinioContainer) Destroy() {
	if err := testcontainers.TerminateContainer(*m.Container); err != nil {
		slog.Error("failed to terminate container", "error", err)
	}
}
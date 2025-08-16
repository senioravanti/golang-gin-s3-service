package handlers

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"mime/multipart"
	"net/http/httptest"
	"net/textproto"

	"testing"

	"github.com/stretchr/testify/assert"

	"senioravanti.ru/internal/s3-service/bootstrap/config"
	"senioravanti.ru/internal/s3-service/rest/handlers"
	"senioravanti.ru/test/helpers"
	"senioravanti.ru/test/testcontainers"
)

// clear; go test -v ./test/rest/handlers
func TestUpload(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions { AddSource: true, Level: slog.LevelDebug, ReplaceAttr: nil }))
	slog.SetDefault(logger)

	minioContainer, err := testcontainers.NewMinioContainer()
	if err != nil {
		slog.Error("failed to start minio container", "error", err)
		return
	}
	defer minioContainer.Destroy()

	assert := assert.New(t)

	s3Config := &config.S3Config{
		AccessKey: minioContainer.Container.Username,
		SecretKey: minioContainer.Container.Password,
		Url: fmt.Sprintf("http://%s", minioContainer.Url),
	}

	s3Client, err := config.NewS3Client(s3Config)
	assert.Nil(err, "s3Client should be successfully created")
	
	bucketHandler := &handlers.BucketHandler{
		S3Client: s3Client,
	}
	objectHandler := &handlers.ObjectHandler{
		S3Client: s3Client,
	}

	output, err := bucketHandler.CreateBucket(context.TODO(), "test")
	assert.Nil(err, "bucket should be successfully created")
	slog.Debug("bucket has been successfully created", "output", output)

	body := &bytes.Buffer{}
	tempFile := helpers.CreateTempTextFile("test content")
	writer := multipart.NewWriter(body)
	
	fileHeader := make(textproto.MIMEHeader)
	fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`,
		filepath.Base(tempFile.Name())))
	fileHeader.Set("Content-Type", "text/plain")

	part, _ := writer.CreatePart(fileHeader)
	tempFileContent, _ := os.ReadFile(tempFile.Name())
	part.Write(tempFileContent)
	writer.Close()

	mockRequest := httptest.NewRequest("POST", "/api/v1/objects/test", body)
	mockRequest.Header.Set("Content-Type", writer.FormDataContentType())

	mockRequestFile, mockRequestFileHeader, _ := mockRequest.FormFile("file")
	defer mockRequestFile.Close()
	objectResponse, err := objectHandler.Upload(
		context.TODO(),
		&handlers.UploadRequest{
			BucketName: "test",
			FileHeader: mockRequestFileHeader,
			File: mockRequestFile,
		},
	)
	assert.Nil(err, "test file should be successfully uploaded")

	slog.Debug("test file has been successfully uploaded !", "objectResponse", objectResponse)
	assert.True(strings.HasSuffix(objectResponse.Key, ".txt"))
	assert.True(objectResponse.Size > int64(0), "file size should be greater then zero")
	assert.Equal("text/plain", objectResponse.ContentType)
}
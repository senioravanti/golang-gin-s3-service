package handlers

import (
	"context"
	"fmt"
	"net/textproto"
	"time"

	"log/slog"

	"mime/multipart"

	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/google/uuid"

	"senioravanti.ru/internal/s3-service/helpers"
	"senioravanti.ru/internal/s3-service/model"
)

type ObjectHandler struct {
	S3Client *s3.Client
}

type UploadRequest struct {
	BucketName string
	File multipart.File
	FileHeader *multipart.FileHeader
}

func getContentType(
	header textproto.MIMEHeader,
) string {
	contentType := header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func (oh *ObjectHandler) Upload(
	ctx context.Context,
	request *UploadRequest,
) (*model.ObjectResponse, error) {
	path := fmt.Sprintf("/%v%s", uuid.New(), filepath.Ext(request.FileHeader.Filename))
	contentType := getContentType(request.FileHeader.Header)

	output, err := oh.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(request.BucketName),
		Key: aws.String(path),
		Body: request.File,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		helpers.ErrWithCause("failed to put object", err)
	}
	slog.Debug("upload result", "objectOutput", output)

	objectResponse := &model.ObjectResponse{
		Key: path,
		Size: request.FileHeader.Size,
		ContentType: contentType,
		LastModified: time.Now().Format(time.RFC3339),
	}
	return objectResponse, nil
}

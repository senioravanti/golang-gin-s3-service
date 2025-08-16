package routers

import (
	"net/http"

	"context"
	"github.com/gin-gonic/gin"

	"senioravanti.ru/internal/s3-service/bootstrap"
	"senioravanti.ru/internal/s3-service/helpers"
	"senioravanti.ru/internal/s3-service/rest/handlers"
)

type ObjectRouter struct {
	App *bootstrap.Application
	RouterGroup *gin.RouterGroup
}

func (or *ObjectRouter) SetUp() {
	objectHandler := &handlers.ObjectHandler{ 
		S3Client: or.App.S3Client,
	}

	or.RouterGroup.POST("/:bucketName", func (c *gin.Context) {
		file, fileHeader, err := c.Request.FormFile("file")
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.Error(helpers.ErrWithCause("failed to get file from request", err))
			return
		}
		defer file.Close()
		
		objectResponse, err := objectHandler.Upload(
			context.TODO(),	
			&handlers.UploadRequest{
				BucketName: c.Param("bucketName"),
				FileHeader: fileHeader,
				File: file,
			},
		)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Error(helpers.ErrWithCause("failed to upload the file", err))
			return
		}

		c.JSON(http.StatusOK, objectResponse)
	})
}

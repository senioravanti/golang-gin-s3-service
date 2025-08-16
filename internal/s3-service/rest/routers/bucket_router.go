package routers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"senioravanti.ru/internal/s3-service/bootstrap"
	"senioravanti.ru/internal/s3-service/helpers"
	"senioravanti.ru/internal/s3-service/rest/handlers"
)

type BucketRouter struct {
	App *bootstrap.Application
	RouterGroup *gin.RouterGroup
}

func (br *BucketRouter) SetUp() {
	bucketHandler := &handlers.BucketHandler{
		S3Client: br.App.S3Client,
	}

	br.RouterGroup.POST("/:bucketName", func (c *gin.Context) {
		bucketName := c.Param("bucketName")
		_, err := bucketHandler.CreateBucket(context.TODO(), bucketName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Error(helpers.ErrWithCause("failed to create the bucket", err))
			return
		}

		c.Status(http.StatusOK)
	})
}

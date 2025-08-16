package routers

import (
	"senioravanti.ru/internal/s3-service/bootstrap"
	"senioravanti.ru/internal/s3-service/rest/middlewares"
)

type Router interface {
	SetUp()
} 

func SetUp(app *bootstrap.Application) {
	app.Gin.Use(
		middlewares.Cors,
		middlewares.HandleError,
	)

	v1 := app.Gin.Group("/api/v1")

	routers := []Router {
		&ObjectRouter{ app, v1.Group("/objects") },
		&BucketRouter{ app, v1.Group("/buckets") },
	}
	
	for _, router := range routers {
		router.SetUp()
	}
}

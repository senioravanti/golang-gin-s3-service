package middlewares

import (
	"log/slog"

	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

type ProblemDetail struct {
	Timestamp string `json:"timestamp"`
	Title string `json:"title"`
	Detail string `json:"detail"`
}

func HandleError(
	c *gin.Context,
) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err
		slog.Error("unhandled error", "error", err)

		problemDetail := &ProblemDetail{
			Timestamp: time.Now().Format(time.RFC3339),
			Title: "INTERNAL SERVER ERROR",
			Detail: "try your request later",
		}

		c.JSON(http.StatusInternalServerError, problemDetail)
		c.Header("Content-Type", "application/problem+json")
	}
}
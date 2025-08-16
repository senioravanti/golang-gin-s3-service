package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "false")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, HEAD, GET, POST")
	if c.Request.Method == "Options" {
		c.AbortWithStatus(204)
		return
	}
	c.Next()
}
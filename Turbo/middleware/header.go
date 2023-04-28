package middleware

import "github.com/gin-gonic/gin"

func ServerHeaderMiddleware(serverName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Server", serverName)
		c.Next()
	}
}

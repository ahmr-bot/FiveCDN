package main

import "github.com/gin-gonic/gin"

func serverHeaderMiddleware(serverName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Server", serverName)
		c.Next()
	}
}

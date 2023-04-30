package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// CORSMiddleware 为Gin绑定CORS跨域访问处理中间件
func CORSMiddleware() gin.HandlerFunc {
	if viper.GetBool("cors.allow_all_origins") != true {
		config := cors.DefaultConfig()
		origins := viper.GetStringSlice("cors.allow_origins")
		methods := viper.GetStringSlice("cors.allow_methods")
		headers := viper.GetStringSlice("cors.allow_headers")
		credentials := viper.GetBool("cors.allow_credentials")
		exposeHeaders := viper.GetStringSlice("cors.expose_headers")

		if len(origins) > 0 {
			config.AllowOrigins = origins
		}
		if len(methods) > 0 {
			config.AllowMethods = methods
		}
		if len(headers) > 0 {
			config.AllowHeaders = headers
		}
		config.AllowCredentials = credentials
		if len(exposeHeaders) > 0 {
			config.ExposeHeaders = exposeHeaders
		}

		// 使用CORS中间件
		return cors.New(config)
	}
	return cors.Default()
}

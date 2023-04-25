package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	config := NewConfig("config.toml")

	// 创建白名单对象并加载初始白名单
	whiteList := NewWhiteList()
	if err := whiteList.LoadFromFile(config.WhiteListURL); err != nil {
		panic(err)
	}

	// 定时更新白名单
	go whiteList.UpdatePeriodically(config.WhiteListURL, config.WhiteListUpdateInterval)

	// 创建 Gin 引擎
	engine := gin.Default()

	// 模式切换
	if viper.GetBool("debug.gin_mode") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 添加自定义响应头
	engine.Use(serverHeaderMiddleware(config.ServerName))

	// 添加白名单验证中间件
	engine.Use(whiteListMiddleware(whiteList))

	// 注册路由
	registerRoutes(engine)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if err := engine.Run(addr); err != nil {
		panic(err)
	}
}

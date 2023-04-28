package main

import (
	"flag"
	"fmt"
	"github.com/ahmr-bot/MECDN/Turbo/middleware"
	"github.com/ahmr-bot/MECDN/Turbo/pkg"
	"github.com/ahmr-bot/MECDN/Turbo/pkg/config"
	"github.com/ahmr-bot/MECDN/Turbo/pkg/router"
	"github.com/gin-gonic/gin"
)

var (
	ConfigPath string
)

func init() {
	flag.StringVar(&ConfigPath, "c", "", "配置文件路径")
	flag.Parse()
}
func main() {
	// 初始化配置
	config := config.NewConfig(ConfigPath)
	// 创建白名单对象并加载初始白名单
	whiteList := middleware.NewWhiteList()
	if err := whiteList.LoadFromFile(config.WhiteListURL); err != nil {
		panic(err)
	}

	// 定时更新白名单
	go whiteList.UpdatePeriodically(config.WhiteListURL, config.WhiteListUpdateInterval)

	// 设置gin模式
	pkg.SetMode()

	// 创建 Gin 引擎
	engine := gin.Default()

	// 添加自定义响应头
	engine.Use(middleware.ServerHeaderMiddleware(config.ServerName))

	// 添加白名单验证中间件
	engine.Use(middleware.WhiteListMiddleware(whiteList))

	// 注册路由
	router.RegisterRoutes(engine)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	if err := engine.Run(addr); err != nil {
		panic(err)
	}
}

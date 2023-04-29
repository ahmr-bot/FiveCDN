package main

import (
	"flag"
	"fmt"
	"github.com/ahmr-bot/MECDN/Turbo/middleware"
	"github.com/ahmr-bot/MECDN/Turbo/pkg"
	"github.com/ahmr-bot/MECDN/Turbo/pkg/config"
	"github.com/ahmr-bot/MECDN/Turbo/pkg/router"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

var (
	ConfigPath string
)

func init() {
	flag.StringVar(&ConfigPath, "c", "config.toml", "配置文件路径")
	flag.Parse()
}
func main() {
	// 初始化配置
	Config := config.NewConfig(ConfigPath)

	// 设置gin模式
	pkg.SetMode()

	// 创建 Gin 引擎
	engine := gin.Default()

	// 添加自定义响应头
	engine.Use(middleware.ServerHeaderMiddleware(Config.PoweredBy))

	// 添加白名单验证中间件
	if viper.GetBool("whitelist.enabled") == true {
		// 创建白名单对象并加载初始白名单
		whiteList := middleware.NewWhiteList()
		if err := whiteList.LoadFromFile(Config.WhiteListURL); err != nil {
			panic(err)
		}

		// 定时更新白名单
		go whiteList.UpdatePeriodically(Config.WhiteListURL, Config.WhiteListUpdateInterval)

		log.Printf("白名单验证已启用")
		engine.Use(middleware.WhiteListMiddleware(whiteList, Config.ServerName))
	} else {
		log.Printf("白名单验证已关闭")
	}

	// 注册路由
	router.RegisterRoutes(engine)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", Config.Host, Config.Port)
	log.Printf("服务器已启动，监听地址：" + Config.Host + ":" + fmt.Sprint(Config.Port))
	if err := engine.Run(addr); err != nil {
		panic(err)
	}
}

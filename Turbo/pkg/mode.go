package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

func SetMode() {
	// 模式切换
	if viper.GetBool("debug.debug") == true {
		gin.SetMode(gin.DebugMode)
		log.Printf("当前模式：Debug")
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Printf("当前模式：Release")
	}
}

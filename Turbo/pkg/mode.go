package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetMode() {
	// 模式切换
	if viper.GetBool("debug.debug") == true {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

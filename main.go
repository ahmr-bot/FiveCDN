package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WhiteList struct {
	sync.RWMutex
	list map[string]bool
}

func NewWhiteList() *WhiteList {
	return &WhiteList{list: make(map[string]bool)}
}

func (wl *WhiteList) LoadFromFile(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			// handle the error
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	lines := strings.Split(string(body), "\n")
	wl.Lock()
	defer wl.Unlock()
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			wl.list[line] = true
		}
	}

	fmt.Println("Loaded whitelist file from", url)

	return nil
}

func (wl *WhiteList) Contains(path string) bool {
	wl.RLock()
	defer wl.RUnlock()
	for k := range wl.list {
		if strings.HasSuffix(k, "/*") {
			prefix := strings.TrimSuffix(k, "/*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		} else if path == k {
			return true
		}
	}
	return false
}
func main() {
	// 初始化配置
	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	// 设置 Gin 的模式
	ginMode := viper.GetString("mode")
	if ginMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	// 创建白名单对象
	whiteList := NewWhiteList()
	// 加载初始白名单
	err := whiteList.LoadFromFile("https://mecdn.mcserverx.com/gh/ahmr-bot/MECDN-WhiteList/master/list.txt")
	if err != nil {
		panic(err)
	}

	// 定时更新白名单
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			err := whiteList.LoadFromFile("https://mecdn.mcserverx.com/gh/ahmr-bot/MECDN-WhiteList/master/list.txt")
			if err != nil {
				fmt.Println("Failed to update whitelist:", err)
			}
		}
	}()

	// 创建 Gin 引擎
	r := gin.Default()

	serverName := viper.GetString("server")
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Server", serverName)
	})

	r.Use(func(c *gin.Context) {
		if whiteList.Contains(c.Request.URL.Path) {
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusForbidden)
	})
	// 遍历所有代理
	for _, p := range viper.Get("proxy").([]interface{}) {
		proxy := p.(map[string]interface{})

		// 注册路由
		r.GET(proxy["path"].(string)+"/*filepath", func(c *gin.Context) {

			// 获取请求路径
			url := c.Request.URL.Path[len(proxy["path"].(string)):]

			//去掉路径末尾的斜杠（如果有）
			if strings.HasSuffix(url, "/") {
				url = url[:len(url)-1]
			}

			// 构建代理 URL
			proxyURL := "https://" + proxy["domain"].(string) + url

			// 发送代理请求
			resp, err := http.Get(proxyURL)
			if err != nil {
				c.AbortWithError(http.StatusBadGateway, err)
				return
			}
			defer resp.Body.Close()

			// 读取响应内容
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// 返回响应内容
			c.Header("Cache-Control", fmt.Sprintf("max-age=%d", viper.GetInt("cache_time")))
			c.Data(http.StatusOK, "", body)
		})
	}

	// 启动服务器
	if err := r.Run(viper.GetString("server.host") + ":" + viper.GetString("server.port")); err != nil {
		panic(err)
	}
}

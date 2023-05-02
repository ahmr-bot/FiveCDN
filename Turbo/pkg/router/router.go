package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

func RegisterRoutes(r *gin.Engine) {
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
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					// handle the error
				}
			}(resp.Body)

			// 获取content-type
			contentType := resp.Header.Get("Content-Type")

			// 读取响应内容
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			// 返回响应内容
			c.Header("Cache-Control", fmt.Sprintf("max-age=%d", viper.GetInt("cache_time")))
			c.Data(http.StatusOK, contentType, body)
		})
	}
}

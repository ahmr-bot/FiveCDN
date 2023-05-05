package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
		if strings.HasSuffix(k, "*") {
			prefix := strings.TrimSuffix(k, "*")
			if strings.HasPrefix(path, prefix) {
				return true
			}
		} else if path == k {
			return true
		}
	}
	return false
}

func (wl *WhiteList) UpdatePeriodically(url string, interval time.Duration) {
	for {
		time.Sleep(interval)
		err := wl.LoadFromFile(url)
		if err != nil {
			fmt.Println("Failed to update whitelist:", err)
		}
	}
}

func WhiteListMiddleware(whiteList *WhiteList, ServerName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if whiteList.Contains(c.Request.URL.Path) {
			c.Next()
			return
		}
		// c.AbortWithStatus(http.StatusForbidden)
		// 在 403页面  给用户返回 The URL is not in the whitelist
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"ip":         c.ClientIP(),
			"message":    "The URL you requested is not in the whitelist",
			"status":     http.StatusForbidden,
			"powered_by": ServerName,
		})

	}
}

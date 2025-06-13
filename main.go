package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	currentClipboard []byte
	mu               sync.RWMutex
)

func main() {
	// 创建 gin router
	router := gin.Default()

	// GET /get_clipboard => 返回剪贴板文本
	router.GET("/get", func(c *gin.Context) {
		mu.RLock()
		data := currentClipboard
		mu.RUnlock()

		c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
	})

	// POST /set => 写入新的剪贴板数据
	router.POST("/set", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil || len(body) == 0 {
			c.String(http.StatusBadRequest, "Invalid body")
			return
		}

		mu.Lock()
		currentClipboard = body
		fmt.Println("Updated:", string(currentClipboard))
		mu.Unlock()

		c.String(http.StatusOK, "Clipboard updated.")
	})

	// 启动服务
	err := router.Run(":8080")
	if err != nil {
		return
	} // 默认监听 localhost:8080
}

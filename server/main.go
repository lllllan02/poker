package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lllllan02/pocker/server"
)

func main() {
	// 创建新的游戏中心
	hub := server.NewHub()
	go hub.Run()

	// 设置 Gin 模式
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// 处理 WebSocket 连接
	r.GET("/ws", func(c *gin.Context) { server.ServeWs(hub, c.Writer, c.Request) })

	// 启动 HTTP 服务器
	r.Run(":8080")
}

package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/lllllan02/pocker/poker"
)

type Hub struct {
	// 游戏
	game *poker.Game

	// 所有连接的客户端
	clients map[string]*Client

	// 广播事件通道
	// 从客户端接受的消息会通过这个通道广播给其他客户端
	broadcast chan BroadcastEvent

	// 注册新客户端的通道
	register chan *Client

	// 注销客户端的通道
	unregister chan *Client
}

// NewHub 创建新的游戏中心
func NewHub() *Hub {
	return &Hub{
		game:       poker.NewGame(),
		clients:    make(map[string]*Client),
		broadcast:  make(chan BroadcastEvent),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 运行游戏中心
func (h *Hub) Run() {
	for {
		select {

		// 注册新客户端
		case client := <-h.register:
			h.clients[client.id] = client

		// 注销客户端
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}

		// 广播事件
		case e := <-h.broadcast:
			for id, client := range h.clients {
				if _, ok := e.ExcludeClients[id]; ok {
					continue
				}

				// 尝试向客户端发送事件
				select {
				// 向客户端发送事件
				case client.send <- e.Event:

				// 如果客户端的发送通道已满或关闭，则删除客户端
				default:
					delete(h.clients, id)
					close(client.send)
				}
			}

		}
	}
}

// ServeWs 处理 WebSocket 连接
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 创建新的客户端
	client := &Client{
		id:       fmt.Sprintf("client_%s", uuid.New().String()[:8]), // 客户端唯一标识
		username: "",                                                // 用户名
		playerId: "",                                                // 玩家唯一标识
		muted:    false,                                             // 是否被禁言
		conn:     conn,                                              // WebSocket 连接
		game:     hub.game,                                          // 游戏
		hub:      hub,                                               // 所属的游戏中心
		send:     make(chan Event),                                  // 发送消息的缓冲通道
	}

	// 注册客户端
	hub.register <- client

	// 启动客户端的读写循环
	go client.readPump()
	go client.writePump()
}

package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lllllan02/pocker/poker"
	"github.com/spf13/cast"
)

// WebSocket 连接的相关常量配置
const (
	// 向对等端写入消息的超时时间
	writeWait = 10 * time.Second

	// 等待对等端 pong 消息的超时时间
	pongWait = 60 * time.Second

	// 发送 ping 消息的间隔时间，必须小于 pongWait
	pingPeriod = (pongWait * 9) / 10

	// 允许接收的最大消息大小
	maxMessageSize = 8192
)

// WebSocket 连接升级器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,                                       // 读缓冲区大小
	WriteBufferSize: 1024,                                       // 写缓冲区大小
	CheckOrigin:     func(r *http.Request) bool { return true }, // 仅用于测试，生产环境需要proper的源检查
}

// Client WebSocket 连接和游戏中心（Hub）之间的中间层
type Client struct {
	id       string          // 客户端唯一标识
	username string          // 用户名
	playerId string          // 玩家唯一标识
	muted    bool            // 是否被禁言
	conn     *websocket.Conn // WebSocket 连接
	game     *poker.Game     // 游戏
	hub      *Hub            // 所属的游戏中心
	send     chan Event      // 发送消息的缓冲通道
}

// readPump 读取 WebSocket 消息并处理
func (c *Client) readPump() {
	// 确保在函数退出时执行清理操作
	defer func() {
		c.disconnectPlayer()  // 断开玩家连接
		c.hub.unregister <- c // 从游戏中心注销客户端
		c.conn.Close()        // 关闭 WebSocket 连接
	}()

	// 设置连接参数
	c.conn.SetReadLimit(maxMessageSize)              // 设置最大消息大小
	c.conn.SetReadDeadline(time.Now().Add(pongWait)) // 设置读取超时
	c.conn.SetPongHandler(func(string) error {       // 设置 pong 消息处理器
		c.conn.SetReadDeadline(time.Now().Add(pongWait)) // 收到 pong 后更新超时时间
		return nil
	})

	// 持续读取消息
	for {
		var event Event
		if err := c.conn.ReadJSON(&event); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 处理事件
		c.processEvent(event)
	}
}

// writePump 将消息写入 WebSocket 连接
func (c *Client) writePump() {
	// 设置定时器，定期发送 ping 消息
	ticker := time.NewTicker(pingPeriod)

	// 确保在函数退出时执行清理操作
	defer func() {
		ticker.Stop()  // 停止定时器
		c.conn.Close() // 关闭 WebSocket 连接
	}()

	// 持续处理发送事件
	for {
		select {

		// 从发送通道接收事件
		case event, ok := <-c.send:
			// 设置写入超时
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			// 如果通道关闭，则关闭连接
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 写入消息
			if err := c.conn.WriteJSON(event); err != nil {
				return
			}

		// 定期发送 ping 消息
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// disconnectPlayer 断开玩家连接
func (c *Client) disconnectPlayer() {}

// processEvent 处理接收的事件
func (c *Client) processEvent(e Event) {
	var err error
	defer func() {
		// 如果处理事件时发生错误，则发送错误事件
		if err != nil {
			c.send <- createErrorEvent(err)
		}
	}()

	switch e.Action {

	// 加入游戏请求
	case EventActionJoin:
		username := cast.ToString(e.Params["username"])
		err = c.handleJoin(username)

	// 入座请求
	case EventActionTakeSeat:
		seatId := cast.ToString(e.Params["seat_id"])
		err = c.handleTakeSeat(seatId)

	// 静音请求
	case EventActionMute:
		muted := cast.ToBool(e.Params["muted"])
		err = c.handleMute(muted)

	// 发送消息
	case EventActionSendMessage:
		username := cast.ToString(e.Params["username"])
		message := cast.ToString(e.Params["message"])
		err = c.handleSendMessage(username, message)

	// 发送信号
	case EventActionSendSignal:
		peerId := cast.ToString(e.Params["peer_id"])
		stream := cast.ToString(e.Params["stream_id"])
		signalData := e.Params["signal_data"]
		err = c.handleSendSignal(peerId, stream, signalData)

	// 处理游戏动作
	default:
		// 如果当前不是玩家回合，则返回错误
		if !c.game.IsPlayerStage() {
			err = fmt.Errorf("you cannot move during the %s stage", c.game.Stage)
			return
		}

		// 如果当前不是玩家回合，则返回错误
		if !c.game.IsPlayerTurn(c.playerId) {
			err = fmt.Errorf("you cannot move out of turn")
			return
		}

		// 处理游戏动作
		switch e.Action {

		// 弃牌
		case EventActionFold:
			err = c.handleFold()

		// 过牌
		case EventActionCheck:
			err = c.handleCheck()

		// 跟注
		case EventActionCall:
			err = c.handleCall()

		// 加注
		case EventActionRaise:
			amount := cast.ToInt(e.Params["amount"])
			err = c.handleRaise(amount)

		default:
			err = fmt.Errorf("invalid action: %s", e.Action)

		}
	}

}

func (c *Client) handleJoin(username string) error {

	return nil
}

func (c *Client) handleTakeSeat(seatId string) error {

	return nil
}

func (c *Client) handleMute(muted bool) error {

	return nil
}

func (c *Client) handleSendMessage(username string, message string) error {

	return nil
}

func (c *Client) handleSendSignal(peerId string, stream string, signalData any) error {

	return nil
}

func (c *Client) handleFold() error {

	return nil
}

func (c *Client) handleCheck() error {

	return nil
}

func (c *Client) handleCall() error {

	return nil
}

func (c *Client) handleRaise(amount int) error {

	return nil
}

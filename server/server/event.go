package server

const (
	// 客户端发给服务端的游戏外动作

	EventActionJoin        = "join"         // 加入请求
	EventActionTakeSeat    = "take_seat"    // 入座请求
	EventActionMute        = "mute"         // 禁言请求
	EventActionSendMessage = "send_message" // 发送消息
	EventActionSendSignal  = "send_signal"  // 发送信号

	// 客户端发给服务端的游戏动作

	EventActionCall  = "call"  // 跟注
	EventActionCheck = "check" // 过牌
	EventActionFold  = "fold"  // 弃牌
	EventActionRaise = "raise" // 加注

	// 服务端发给客户端

	EventActionError = "error" // 错误事件

)

type Event struct {
	Action string         `json:"action"` // 事件类型
	Params map[string]any `json:"params"` // 事件参数
}

type BroadcastEvent struct {
	Event          Event           // 要广播的事件
	ExcludeClients map[string]bool // 排除的客户端列表
}

func createErrorEvent(err error) Event {
	return Event{
		Action: EventActionError,
		Params: map[string]any{
			"error": err.Error(),
		},
	}
}

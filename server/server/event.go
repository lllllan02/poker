package server

type Event struct {
	Action string         `json:"action"` // 事件类型
	Params map[string]any `json:"params"` // 事件参数
}

type BroadcastEvent struct {
	Event          Event           // 要广播的事件
	ExcludeClients map[string]bool // 排除的客户端列表
}

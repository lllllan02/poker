package server

import "github.com/lllllan02/pocker/poker"

// GetActions 获取当前玩家可以执行的动作
func GetActions(g *poker.Game) []string {
	var actions []string

	// 如果玩家可以弃牌，则添加弃牌动作
	if g.CurrentSeat.Player.CanFold(g.BettingRound) {
		actions = append(actions, EventActionFold)
	}

	// 如果玩家可以过牌，则添加过牌动作
	if g.CurrentSeat.Player.CanCheck(g.BettingRound) {
		actions = append(actions, EventActionCheck)
	}

	// 如果玩家可以跟注，则添加跟注动作
	if g.CurrentSeat.Player.CanCall(g.BettingRound) {
		actions = append(actions, EventActionCall)
	}

	// 如果玩家可以加注，则添加加注动作
	if g.CurrentSeat.Player.CanRaise(g.BettingRound) {
		actions = append(actions, EventActionRaise)
	}

	return actions
}

// createClientPlayerMap 创建客户端玩家映射
func createClientPlayerMap(clients map[string]*Client) map[string]string {
	clientPlayerMap := make(map[string]string)
	for _, c := range clients {
		clientPlayerMap[c.id] = c.playerId
	}
	return clientPlayerMap
}

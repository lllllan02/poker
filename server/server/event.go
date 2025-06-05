package server

import (
	"github.com/google/uuid"
	"github.com/lllllan02/pocker/poker"
)

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

	EventActionError      = "error"       // 错误事件
	EventActionOnJoin     = "on_join"     // 加入成功
	EventActionNewMessage = "new_message" // 新消息
	EventActionUpdateGame = "update_game" // 更新游戏
)

type Event struct {
	Action string         `json:"action"` // 事件类型
	Params map[string]any `json:"params"` // 事件参数
}

// 创建错误事件
func createErrorEvent(err error) Event {
	return Event{
		Action: EventActionError,
		Params: map[string]any{
			"error": err.Error(),
		},
	}
}

// 创建加入成功事件
func createOnJoinEvent(userId, username string) Event {
	return Event{
		Action: EventActionOnJoin,
		Params: map[string]any{
			"user_id":  userId,
			"username": username,
		},
	}
}

// 创建新消息事件
func createNewMessageEvent(username, message string) Event {
	return Event{
		Action: EventActionNewMessage,
		Params: map[string]any{
			"id":       uuid.New().String(),
			"username": username,
			"message":  message,
		},
	}
}

func createUpdateGameEvent(c *Client, showCards bool) Event {
	game := c.game
	seats := game.Table.Seats
	players := make([]map[string]any, 0)
	var actionBar map[string]any

	if game.Stage == poker.GameStageWaiting {
		for i := 0; i < seats.Len(); i++ {
			players = append(players, map[string]any{
				"id":         seats.Player.Id,              // 玩家 id
				"name":       seats.Player.Name,            // 玩家名称
				"status":     seats.Player.Status.String(), // 玩家状态
				"is_active":  false,                        // 是否活跃
				"is_dealer":  false,                        // 是否庄家
				"chips":      seats.Player.Chips,           // 筹码
				"has_folded": seats.Player.HasFolded,       // 是否弃牌
				"hole_cards": [2]*poker.Card{},             // 手牌
			})
			seats = seats.Next()
		}

		actionBar = map[string]any{
			"actions":        []string{}, // 动作
			"callAmount":     0,          // 跟注金额
			"chipsInPot":     0,          // 奖池内筹码
			"maxRaiseAmount": 0,          // 最大加注金额
			"minBetAmount":   0,          // 最小下注金额
			"minRaiseAmount": 0,          // 最小加注金额
			"totalChips":     0,          // 总筹码
		}
	} else {
		activePlayer := game.CurrentSeat.Player
		for i := 0; i < seats.Len(); i++ {
			for i := 0; i < seats.Len(); i++ {
				holeCards := [2]*poker.Card{}
				if showCards && !seats.Player.HasFolded {
					holeCards = seats.Player.HoleCards
				}

				players = append(players, map[string]interface{}{
					"id":         seats.Player.Id,
					"name":       seats.Player.Name,
					"status":     seats.Player.Status.String(),
					"isActive":   seats.Player.Id == activePlayer.Id,
					"isDealer":   seats.Player.Id == game.Table.Dealer.Player.Id,
					"chips":      seats.Player.Chips,
					"chipsInPot": game.BettingRound.Bets[seats.Player],
					"hasFolded":  seats.Player.HasFolded,
					"holeCards":  holeCards,
				})

				seats = seats.Next()
			}
		}

		// 如果玩家筹码不足，则将最大加注金额设置为玩家剩余筹码
		callRemainingAmount := game.BettingRound.CallAmount - game.BettingRound.Bets[activePlayer]
		maxRaiseAmount := activePlayer.Chips - callRemainingAmount
		if maxRaiseAmount < 0 {
			maxRaiseAmount = activePlayer.Chips
		}

		// 如果最小加注金额小于最大加注金额，则使用最大加注金额作为最小加注金额
		// 这意味着玩家必须全押才能加注
		minRaiseAmount := game.BettingRound.RaiseByAmount
		if minRaiseAmount > maxRaiseAmount {
			minRaiseAmount = maxRaiseAmount
		}

		actionBar = map[string]interface{}{
			"actions":        GetActions(game),
			"callAmount":     game.BettingRound.CallAmount,
			"chipsInPot":     game.BettingRound.Bets[activePlayer],
			"maxRaiseAmount": maxRaiseAmount,
			"minBetAmount":   game.Table.MinBet,
			"minRaiseAmount": minRaiseAmount,
			"seatID":         activePlayer.Id,
			"totalChips":     activePlayer.Chips,
		}
	}

	table := map[string]interface{}{
		"flop":  game.Table.Flop,
		"pot":   game.Table.Pot.GetTotal(),
		"river": game.Table.River,
		"turn":  game.Table.Turn,
	}

	return Event{
		Action: EventActionUpdateGame,
		Params: map[string]any{
			"action_bar":        actionBar,
			"client_player_map": createClientPlayerMap(c.hub.clients),
			"players":           players,
			"stage":             game.Stage.String(),
			"table":             table,
		},
	}
}

type BroadcastEvent struct {
	Event          Event           // 要广播的事件
	ExcludeClients map[string]bool // 排除的客户端列表
}

// NewBroadcastEvent 创建广播事件
func NewBroadcastEvent(e Event) BroadcastEvent {
	return BroadcastEvent{
		Event:          e,
		ExcludeClients: make(map[string]bool),
	}
}

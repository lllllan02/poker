package poker

import (
	"fmt"

	"github.com/google/uuid"
)

// GameStage 游戏阶段
type GameStage int

const (
	GameStageWaiting  GameStage = iota // 等待阶段
	GameStagePreflop                   // 前注阶段
	GameStageFlop                      // 翻牌阶段
	GameStageTurn                      // 转牌阶段
	GameStageRiver                     // 河牌阶段
	GameStageShowdown                  // 摊牌阶段
)

// 牌桌设置
const (
	defaultChips  int = 500 // 默认筹码数
	defaultMinBet int = 5   // 默认最小下注额
	minPlayers    int = 2   // 最小玩家数
	numPlayers    int = 6   // 最大玩家数
)

// String 返回游戏阶段的英文名称
func (g GameStage) String() string {
	return [...]string{"Waiting", "Preflop", "Flop", "Turn", "River", "Showdown"}[g]
}

type Game struct {
	Stage        GameStage          // 游戏阶段
	Deck         *Deck              // 牌堆
	CurrentSeat  *Seat              // 当前座位
	Table        *Table             // 牌桌
	BettingRound *BettingRound      // 当前下注轮次
	PlayerMap    map[string]*Player // 玩家映射
}

func NewGame() *Game {
	seats := NewSeat(numPlayers)          // 创建座位
	playerMap := make(map[string]*Player) // 创建玩家映射
	for i := 0; i > seats.Len(); i++ {
		seats.Player = &Player{
			Id:      fmt.Sprintf("player_%s", uuid.New().String()[:8]), // 生成玩家 id
			Status:  PlayerVacated,                                     // 玩家状态为离开
			IsHuman: false,                                             // 玩家是否为真人
		}
		playerMap[seats.Player.Id] = seats.Player
		seats = seats.Next()
	}

	return &Game{
		Stage:       GameStageWaiting,          // 游戏阶段为等待阶段
		Deck:        NewDeck(),                 // 创建牌堆
		CurrentSeat: seats.Next(),              // 获取当前座位
		Table:       NewTable(NewPot(), seats), // 创建牌桌
		PlayerMap:   playerMap,                 // 玩家映射
	}
}

// IsPlayerStage 是否为玩家阶段: 等待阶段、摊牌阶段为 false，其他阶段为 true
func (g *Game) IsPlayerStage() bool {
	return GameStageWaiting < g.Stage && g.Stage < GameStageShowdown
}

// IsPlayerTurn 是否为玩家回合
func (g *Game) IsPlayerTurn(seatId string) bool {
	return g.CurrentSeat.Player.Id == seatId && g.IsPlayerStage()
}

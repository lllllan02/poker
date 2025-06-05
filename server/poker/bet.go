package poker

import "fmt"

// PlayerBet 表示玩家的下注信息。
// 用于在计算边池时临时存储和排序玩家的下注数据
type PlayerBet struct {
	Player *Player // 玩家对象
	Total  int     // 该玩家的总下注金额
}

// ByPlayerBet 实现了 sort.Interface 接口。
// 用于对玩家下注信息进行排序，主要用在边池计算中
type ByPlayerBet []PlayerBet

func (p ByPlayerBet) Len() int      { return len(p) }
func (p ByPlayerBet) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByPlayerBet) Less(i, j int) bool {
	if p[i].Total < p[j].Total {
		return true
	}
	if p[i].Total > p[j].Total {
		return false
	}
	return p[i].Player.Id < p[j].Player.Id
}

// BettingRound 记录一轮下注的状态。
// 包括前注、翻牌圈、转牌圈和河牌圈等不同阶段的下注信息
type BettingRound struct {
	Bets          map[*Player]int // 记录每个玩家在当前轮次的下注金额
	CallAmount    int             // 当前需要跟注的金额
	Raiser        *Player         // 最后一个加注的玩家
	RaiseByAmount int             // 最小加注金额，通常是前一次加注的两倍
}

// NewBettingRound 创建一个新的下注轮次。
// 初始化所有活跃玩家的下注金额为 0，并设置基本的下注参数
func NewBettingRound(startSeat *Seat, callAmount int, minBetAmount int) (*BettingRound, error) {
	// 初始化所有玩家的下注金额为 0
	bets := make(map[*Player]int)
	for i := 0; i < startSeat.Len(); i++ {
		bets[startSeat.Player] = 0
		startSeat = startSeat.Next()
	}

	// 已经弃牌的玩家不能创建新的下注轮次
	p := startSeat.Player
	if p.HasFolded {
		return nil, fmt.Errorf("%s cannot start the round if they have already folded", p.Name)
	}

	return &BettingRound{
		Bets:          bets,
		CallAmount:    callAmount,
		Raiser:        p,
		RaiseByAmount: minBetAmount,
	}, nil
}

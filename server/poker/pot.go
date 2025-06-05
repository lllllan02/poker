package poker

// Pot 表示当前游戏中的奖池。
// 记录了所有玩家的下注情况，并支持边池的计算
type Pot struct {
	Bets map[*Player]int // 记录每个玩家的下注金额，key 为玩家指针，value 为下注金额
}

// SidePot 表示边池。
// 当有玩家全下时会产生边池，确保每个玩家只能赢取其下注额度内的奖金
type SidePot struct {
	Players []*Player // 参与该边池的玩家列表
	Total   int       // 边池中的总金额
	MaxBet  int       // 该边池中每个玩家的最大下注额
}

// GetTotal 计算奖池的总金额
func (p *Pot) GetTotal() int {
	total := 0
	for _, betAmount := range p.Bets {
		total += betAmount
	}
	return total
}

// NewPot 创建一个新的奖池
func NewPot() *Pot {
	return &Pot{Bets: make(map[*Player]int)}
}

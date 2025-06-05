package poker

import (
	"errors"
	"fmt"
)

// PlayerStatus 玩家状态的枚举类型
type PlayerStatus int

const (
	PlayerVacated    PlayerStatus = iota // 已离开
	PlayerSittingOut                     // 暂时离座
	PlayerActive                         // 正在游戏
)

// String 返回玩家状态的英文名称
func (p PlayerStatus) String() string {
	return [...]string{"vacated", "sitting-out", "active"}[p]
}

// Player 表示扑克游戏中的一个玩家。
// 它记录了玩家在一轮游戏中的状态，包括从前注到河牌的所有下注轮次
type Player struct {
	Chips     int          // 玩家持有的筹码数
	HasFolded bool         // 是否已弃牌
	HoleCards [2]*Card     // 玩家的两张底牌
	Id        string       // 玩家唯一标识(座位号)
	Name      string       // 玩家名称
	IsHuman   bool         // 是否是人类玩家
	Status    PlayerStatus // 玩家当前状态
}

// PrintHoleCards 打印玩家的两张底牌
func (p *Player) PrintHoleCards() (string, error) {
	if p.HoleCards[0] == nil || p.HoleCards[1] == nil {
		return "", errors.New("the player does not have any hole cards yet")
	}

	hand := fmt.Sprintf("%s %s", p.HoleCards[0].Symbol(), p.HoleCards[1].Symbol())
	return hand, nil
}

// CanFold 检查玩家是否可以弃牌。
func (p *Player) CanFold(b *BettingRound) bool {
	return (p.Status == PlayerActive && // 玩家处于活跃状态
		!p.HasFolded && // 尚未弃牌
		p.Chips > 0 && // 还有筹码
		b.Bets[p] < b.CallAmount) // 当前下注小于跟注金额
}

// Fold 执行弃牌操作。
func (p *Player) Fold(b *BettingRound) error {
	// 玩家处于活跃状态
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}

	// 玩家还有筹码
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}

	// 玩家尚未弃牌
	if p.HasFolded {
		return fmt.Errorf("%s has already folded", p.Name)
	}

	// 玩家当前下注小于跟注金额
	if b.Bets[p] >= b.CallAmount {
		return fmt.Errorf("%s has matched or exceeded the current bets", p.Name)
	}

	p.HasFolded = true
	return nil
}

// CanCheck 检查玩家是否可以过牌。
func (p *Player) CanCheck(b *BettingRound) bool {
	return (p.Status == PlayerActive && // 玩家处于活跃状态
		!p.HasFolded && // 尚未弃牌
		p.Chips > 0 && // 还有筹码
		b.Bets[p] == b.CallAmount) // 当前下注等于跟注金额
}

// Check 执行过牌操作，将行动权传给下一位玩家。
func (p *Player) Check(b *BettingRound) error {
	// 玩家处于活跃状态
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}

	// 玩家没有筹码
	if p.Chips == 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}

	// 玩家已经弃牌
	if p.HasFolded {
		return fmt.Errorf("%s can't check when folded", p.Name)
	}

	// 玩家当前下注小于跟注金额
	if b.Bets[p] < b.CallAmount {
		return fmt.Errorf("%s can't check when someone has bet", p.Name)
	}
	return nil
}

// CanCall 检查玩家是否可以跟注。
func (p *Player) CanCall(b *BettingRound) bool {
	return (p.Status == PlayerActive && // 玩家处于活跃状态
		!p.HasFolded && // 尚未弃牌
		p.Chips > 0 && // 还有筹码
		b.Bets[p] < b.CallAmount) // 当前下注小于跟注金额
}

// Call 执行跟注操作。
func (p *Player) Call(t *Table, b *BettingRound) error {
	// 玩家处于活跃状态
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}

	// 玩家没有筹码
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}

	// 玩家已经弃牌
	if p.HasFolded {
		return fmt.Errorf("%s can't call when folded", p.Name)
	}

	// 如果玩家已下注的筹码超过跟注金额，则不能跟注
	chipsInPot := b.Bets[p]
	if chipsInPot >= b.CallAmount {
		return fmt.Errorf("%s has more chips wagered (%d) than the call amount (%d)", p.Name, chipsInPot, b.CallAmount)
	}

	callAmount := b.CallAmount - chipsInPot

	// 如果玩家的筹码不足以跟注，则全下
	if callAmount > p.Chips {
		callAmount = p.Chips
	}

	t.Pot.Bets[p] += callAmount
	b.Bets[p] += callAmount
	p.Chips -= callAmount

	return nil
}

// CanRaise 检查玩家是否可以加注。
//
// 如果玩家的筹码不足以满足最小加注额，则可以选择全下
func (p *Player) CanRaise(b *BettingRound) bool {
	return (p.Status == PlayerActive && // 玩家处于活跃状态
		!p.HasFolded && // 尚未弃牌
		p.Chips > 0 && // 还有筹码
		p.Chips >= b.CallAmount-b.Bets[p]) // 剩余筹码足够跟注
}

// Raise 执行加注操作。
func (p *Player) Raise(t *Table, b *BettingRound, raiseAmount int) error {
	// 根据当前是否有人下注，决定操作名称
	actionLabel := "bet"
	if b.CallAmount > 0 {
		actionLabel = "raise"
	}

	// 玩家处于活跃状态
	if p.Status != PlayerActive {
		return fmt.Errorf("%s is not in the hand", p.Name)
	}

	// 玩家没有筹码
	if p.Chips <= 0 {
		return fmt.Errorf("%s can't play with no chips", p.Name)
	}

	// 玩家已经弃牌
	if p.HasFolded {
		return fmt.Errorf("%s can't %s when folded", p.Name, actionLabel)
	}

	// 计算需要补充的筹码数
	chipsInPot := b.Bets[p]
	chipsNeeded := raiseAmount - chipsInPot

	// 检查玩家是否有足够的筹码
	if chipsNeeded > p.Chips {
		return fmt.Errorf("%s does not have enough chips (%d) to %s (%d)", p.Name, p.Chips, actionLabel, chipsNeeded)
	}

	// 计算最小加注额
	minRaiseTo := b.CallAmount + b.RaiseByAmount

	// 如果玩家有足够的筹码满足最小加注，则必须至少加注最小额度
	if raiseAmount < minRaiseTo && minRaiseTo-chipsInPot < p.Chips {
		return fmt.Errorf("%s's raise (%d) is less than the minimum %s (%d)", p.Name, raiseAmount, actionLabel, minRaiseTo)
	}

	// 只有当加注金额大于等于最小加注时，才更新最小加注额
	if raiseAmount >= minRaiseTo {
		b.RaiseByAmount = raiseAmount - b.CallAmount
	}

	// 更新游戏状态
	b.CallAmount = raiseAmount
	b.Raiser = p
	t.Pot.Bets[p] += chipsNeeded
	b.Bets[p] += chipsNeeded
	p.Chips -= chipsNeeded

	return nil
}

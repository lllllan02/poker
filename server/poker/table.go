package poker

import "fmt"

// Table 表示当前扑克游戏的牌桌状态。
// 包含了游戏进行所需的所有关键信息，如座位分布、庄家位置、下注情况和公共牌等
type Table struct {
	Seats      *Seat    // 所有座位，使用循环链表结构连接，从第一个座位开始
	Dealer     *Seat    // 庄家座位，每轮游戏结束后按顺时针移动
	SmallBlind *Seat    // 小盲注座位，位于庄家的下一个位置
	BigBlind   *Seat    // 大盲注座位，位于小盲注的下一个位置
	MinBet     int      // 最小下注额，通常等于大盲注的金额
	Pot        *Pot     // 当前奖池，记录所有玩家的下注金额
	Flop       [3]*Card // 公共牌：翻牌，游戏中首先发出的三张公共牌
	Turn       *Card    // 公共牌：转牌，第四张公共牌
	River      *Card    // 公共牌：河牌，第五张也是最后一张公共牌
}

// DealHands 为所有玩家发放手牌
func (t *Table) DealHands(d *Deck) {
	activePlayers := t.Seats.GetActivePlayers()
	hands := make([][2]*Card, len(activePlayers))
	for rounds := range []int{0, 1} {
		for i := range hands {
			card, _ := d.GetNextCard()
			hands[i][rounds] = card
		}
	}

	for i := range hands {
		activePlayers[i].HoleCards = hands[i]
	}
}

// TakeSmallBlind 收取小盲注
func (t *Table) TakeSmallBlind(b *BettingRound) error {
	p := t.SmallBlind.Player
	if p.Chips < t.MinBet {
		return fmt.Errorf("%s does not have enough chips to play", p.Name)
	}

	smallBlind := t.MinBet
	p.Chips -= smallBlind
	t.Pot.Bets[p] += smallBlind
	b.Bets[p] = smallBlind

	return nil
}

// TakeBigBlind 收取大盲注
func (t *Table) TakeBigBlind(b *BettingRound) error {
	p := t.BigBlind.Player
	if p.Chips < t.MinBet {
		return fmt.Errorf("%s does not have enough chips to play", p.Name)
	}

	bigBlind := t.MinBet * 2
	p.Chips -= bigBlind
	t.Pot.Bets[p] += bigBlind
	b.Bets[p] = bigBlind
	b.CallAmount = bigBlind
	b.RaiseByAmount = bigBlind
	return nil
}

// DealFlop 发放三张公共牌（翻牌）
func (t *Table) DealFlop(d *Deck) {
	for i := range t.Flop {
		card, _ := d.GetNextCard()
		t.Flop[i] = card
	}
}

// DealTurn 发放第四张公共牌（转牌）
func (t *Table) DealTurn(d *Deck) {
	card, _ := d.GetNextCard()
	t.Turn = card
}

// DealRiver 发放第五张公共牌（河牌）
func (t *Table) DealRiver(d *Deck) {
	card, _ := d.GetNextCard()
	t.River = card
}

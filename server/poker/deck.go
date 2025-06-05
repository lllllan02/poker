package poker

import (
	"fmt"
	"math/rand"
)

// DeckSize 一副扑克牌的总数量
const DeckSize = 52

// Deck 表示一副扑克牌
type Deck struct {
	Cards            []Card // 所有扑克牌的切片
	CurrentCardIndex int    // 当前发牌位置的索引
}

// GetNextCard 从牌堆中获取下一张牌
func (d *Deck) GetNextCard() (*Card, error) {
	if d.CurrentCardIndex >= DeckSize {
		return nil, fmt.Errorf("no more cards left in deck")
	}

	card := d.Cards[d.CurrentCardIndex]
	d.CurrentCardIndex++
	return &card, nil
}

// NewDeck 创建一副新的扑克牌
func NewDeck() Deck {
	suits := []CardSuit{Clubs, Diamonds, Hearts, Spades}
	ranks := []CardRank{Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace}
	cards := make([]Card, DeckSize)

	// 创建一副完整的扑克牌，包含每种花色和点数的组合
	i := 0
	for _, s := range suits {
		for _, r := range ranks {
			cards[i] = Card{
				Rank: r,
				Suit: s,
			}
			i++
		}
	}

	// 使用随机算法打乱牌堆
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })

	return Deck{Cards: cards, CurrentCardIndex: 0}
}

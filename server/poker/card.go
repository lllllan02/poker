package poker

import "fmt"

// CardSuit 扑克牌花色的枚举类型
type CardSuit int

const (
	Clubs    CardSuit = iota // 梅花
	Diamonds                 // 方块
	Hearts                   // 红心
	Spades                   // 黑桃
)

// String 返回花色的英文名称
func (s CardSuit) String() string {
	return [...]string{"Clubs", "Diamonds", "Hearts", "Spades"}[s]
}

// Symbol 返回花色的符号表示
func (s CardSuit) Symbol() string {
	return [...]string{"♣", "♦", "♥", "♠"}[s]
}

// CardRank 扑克牌点数的枚举类型
type CardRank int

const (
	Two   CardRank = iota // 2
	Three                 // 3
	Four                  // 4
	Five                  // 5
	Six                   // 6
	Seven                 // 7
	Eight                 // 8
	Nine                  // 9
	Ten                   // 10
	Jack                  // J
	Queen                 // Q
	King                  // K
	Ace                   // A
)

// String 返回点数的英文名称
func (r CardRank) String() string {
	return [...]string{
		"Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King", "Ace",
	}[r]
}

// Symbol 返回点数的符号表示
func (r CardRank) Symbol() string {
	return [...]string{
		"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A",
	}[r]
}

// Card 表示一张扑克牌，包含点数和花色
type Card struct {
	Rank CardRank `json:"rank"` // 点数
	Suit CardSuit `json:"suit"` // 花色
}

// String 返回扑克牌的完整描述，例如："Ace of Spades"
func (c *Card) String() string {
	return fmt.Sprintf("%s of %s", c.Rank, c.Suit)
}

// Symbol 返回扑克牌的简短符号表示，例如："A♠"
func (c *Card) Symbol() string {
	return fmt.Sprintf("%s%s", c.Rank.Symbol(), c.Suit.Symbol())
}

// ByCard 实现了扑克牌的排序接口，按点数升序排列
// 注意：由于没有使用花色作为次要排序条件，因此排序结果可能不是唯一的
type ByCard []Card

func (c ByCard) Len() int           { return len(c) }
func (c ByCard) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ByCard) Less(i, j int) bool { return c[i].Rank < c[j].Rank }

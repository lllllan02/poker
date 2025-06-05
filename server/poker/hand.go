package poker

import "sort"

// Comparison 用于比较手牌大小的枚举类型 (>, =, <)
type Comparison int

const (
	LessThan    Comparison = -1 // 小于
	EqualTo     Comparison = 0  // 等于
	GreaterThan Comparison = 1  // 大于
)

// HandRank 扑克牌牌型的等级枚举
type HandRank int

const (
	HighCard      HandRank = iota // 高牌
	OnePair                       // 一对
	TwoPair                       // 两对
	ThreeOfAKind                  // 三条
	Straight                      // 顺子
	Flush                         // 同花
	FullHouse                     // 葫芦
	FourOfAKind                   // 四条
	StraightFlush                 // 同花顺
	RoyalFlush                    // 皇家同花顺
)

// String 返回手牌牌型的英文名称
func (h HandRank) String() string {
	return [...]string{"High Card", "One Pair", "Two Pair", "Three Of A Kind", "Straight", "Flush", "Full House", "Four Of A Kind", "Straight Flush", "Royal Flush"}[h]
}

// Hand 表示玩家手中的牌型。
// 包括牌型等级和用于平局判定的牌点数列表
type Hand struct {
	Rank        HandRank   // 牌型等级
	TieBreakers []CardRank // 用于平局判定的牌点数列表
}

// CompareLess 比较当前手牌是否小于另一副手牌
func (h *Hand) CompareLess(otherHand *Hand) bool {
	if h.Rank < otherHand.Rank {
		return true
	}

	if h.Rank > otherHand.Rank {
		return false
	}

	// 只有在相同牌型的情况下才比较平局判定值
	for i := range h.TieBreakers {
		if h.TieBreakers[i] < otherHand.TieBreakers[i] {
			return true
		}
		if h.TieBreakers[i] > otherHand.TieBreakers[i] {
			return false
		}
	}

	return false
}

// ByHand 实现了手牌的排序接口
type ByHand []Hand

func (h ByHand) Len() int           { return len(h) }
func (h ByHand) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h ByHand) Less(i, j int) bool { return h[i].CompareLess(&h[j]) }

// IsHand 是一个函数类型，用于检查特定的牌型
type IsHand func(cs [5]Card) *Hand

// IsRoyalFlush 检查是否是皇家同花顺。
// 皇家同花顺是同一花色的 10-J-Q-K-A
func IsRoyalFlush(cs [5]Card) *Hand {
	straightFlush := IsStraightFlush(cs)

	if straightFlush == nil {
		return nil
	}

	if straightFlush.TieBreakers[0] != Ace {
		return nil
	}

	return &Hand{
		Rank:        RoyalFlush,
		TieBreakers: []CardRank{},
	}
}

// IsStraightFlush 检查是否是同花顺。
// 同花顺是同一花色的连续五张牌
func IsStraightFlush(cs [5]Card) *Hand {
	flush := IsFlush(cs)
	straight := IsStraight(cs)

	if flush == nil || straight == nil {
		return nil
	}

	return &Hand{
		Rank:        StraightFlush,
		TieBreakers: straight.TieBreakers,
	}
}

// IsFourOfAKind 检查是否是四条。
// 四条是四张相同点数的牌加上一张单牌
func IsFourOfAKind(cs [5]Card) *Hand {
	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	isHand := false     // 是否是四条
	handStrength := Two // 四条的点数
	kicker := Two       // 单牌的点数，用于平局判定
	for rank, count := range rankCount {
		if count == 4 {
			isHand = true
			handStrength = rank
		}

		if count == 1 {
			kicker = rank
		}
	}

	if !isHand {
		return nil
	}

	return &Hand{
		Rank:        FourOfAKind,
		TieBreakers: []CardRank{handStrength, kicker},
	}
}

// IsFullHouse 检查是否是葫芦。
// 葫芦是三张相同点数的牌加上一对
func IsFullHouse(cs [5]Card) *Hand {
	threeOfAKind := IsThreeOfAKind(cs)
	pair := IsOnePair(cs)
	if threeOfAKind == nil || pair == nil {
		return nil
	}

	handStrength := threeOfAKind.TieBreakers[0]
	kicker := pair.TieBreakers[0]

	return &Hand{
		Rank:        FullHouse,
		TieBreakers: []CardRank{handStrength, kicker},
	}
}

// IsFlush 检查是否是同花。
// 同花是五张相同花色的牌
func IsFlush(cs [5]Card) *Hand {
	suit := cs[0].Suit
	for _, c := range cs {
		if c.Suit != suit {
			return nil
		}
	}

	sort.Sort(sort.Reverse(ByCard(cs[:])))
	tieBreakers := make([]CardRank, len(cs))
	for i := range cs {
		tieBreakers[i] = cs[i].Rank
	}

	return &Hand{
		Rank:        Flush,
		TieBreakers: tieBreakers,
	}
}

// IsStraight 检查是否是顺子。
// 顺子是五张连续点数的牌
func IsStraight(cs [5]Card) *Hand {
	sort.Sort(ByCard(cs[:]))

	// 处理 Ace-5 的边缘情况
	if cs[0].Rank == Two && cs[1].Rank == Three && cs[2].Rank == Four && cs[3].Rank == Five && cs[4].Rank == Ace {
		return &Hand{
			Rank:        Straight,
			TieBreakers: []CardRank{Five},
		}
	}

	for i := 0; i < len(cs)-1; i++ {
		if cs[i].Rank+1 != cs[i+1].Rank {
			return nil
		}
	}

	handStrength := cs[4].Rank
	return &Hand{
		Rank:        Straight,
		TieBreakers: []CardRank{handStrength},
	}
}

// IsThreeOfAKind 检查是否是三条。
// 三条是三张相同点数的牌加上两张单牌
func IsThreeOfAKind(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	isHand := false // 是否三条
	tieBreakers := make([]CardRank, 0)
	for rank, count := range rankCount {
		if count == 3 {
			isHand = true
			tieBreakers = append(tieBreakers, rank)
		}
	}

	if !isHand {
		return nil
	}

	// 添加单牌作为平局判定值
	for _, c := range cs {
		if rankCount[c.Rank] != 3 {
			tieBreakers = append(tieBreakers, c.Rank)
		}
	}

	return &Hand{
		Rank:        ThreeOfAKind,
		TieBreakers: tieBreakers,
	}
}

// IsTwoPair 检查是否是两对。
// 两对是两个对子加上一张单牌
func IsTwoPair(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	pairs := 0 // 对子的数量
	tieBreakers := make([]CardRank, 0)
	for rank, count := range rankCount {
		if count == 2 {
			pairs++
			tieBreakers = append(tieBreakers, rank)
		}
	}

	if pairs != 2 {
		return nil
	}

	// 添加单牌作为平局判定值
	for _, c := range cs {
		if rankCount[c.Rank] != 2 {
			tieBreakers = append(tieBreakers, c.Rank)
		}
	}

	return &Hand{
		Rank:        TwoPair,
		TieBreakers: tieBreakers,
	}
}

// IsOnePair 检查是否是一对。
// 一对是一个对子加上三张单牌
func IsOnePair(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	rankCount := make(map[CardRank]int)
	for _, c := range cs {
		rankCount[c.Rank]++
	}

	pairs := 0 // 对子的数量
	tieBreakers := make([]CardRank, 0)
	for rank, count := range rankCount {
		if count == 2 {
			pairs++
			tieBreakers = append(tieBreakers, rank)
		}
	}

	if pairs != 1 {
		return nil
	}

	for _, c := range cs {
		if rankCount[c.Rank] != 2 {
			tieBreakers = append(tieBreakers, c.Rank)
		}
	}

	return &Hand{
		Rank:        OnePair,
		TieBreakers: tieBreakers,
	}
}

// IsHighCard 检查是否为高牌。
// 当不符合任何其他牌型时，使用最大点数的牌作为高牌
func IsHighCard(cs [5]Card) *Hand {
	sort.Sort(sort.Reverse(ByCard(cs[:])))

	tieBreakers := make([]CardRank, len(cs))
	for i := range cs {
		tieBreakers[i] = cs[i].Rank
	}

	return &Hand{
		Rank:        HighCard,
		TieBreakers: tieBreakers,
	}
}

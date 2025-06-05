package poker

import "sort"

// PlayerHand 表示玩家在当前可能组合中的最佳手牌
type PlayerHand struct {
	ChipsWon int     // 赢得的筹码数
	Hand     *Hand   // 最佳手牌组合
	Player   *Player // 玩家
}

// FindWinningHands 找出拥有最佳手牌的玩家。
func FindWinningHands(players []*Player, t *Table) []PlayerHand {
	winners := make([]PlayerHand, 0)

	// 找出拥有相同牌型的玩家
	for i := range players {
		bestHand := GetBestHand(players[i], t)

		if len(winners) == 0 {
			// 如果还没有赢家，将当前玩家设为默认赢家
			winners = append(winners, PlayerHand{Hand: bestHand, Player: players[i]})
		} else {
			// 比较当前玩家的最佳手牌与已有的最佳手牌
			result := CompareHand(bestHand, winners[0].Hand)
			if result == GreaterThan {
				// 如果找到更好的手牌，将当前玩家设为唯一赢家
				winners = []PlayerHand{
					{Hand: bestHand, Player: players[i]},
				}
			} else if result == EqualTo {
				// 如果是平局，玩家将共享奖池
				winners = append(winners, PlayerHand{Hand: bestHand, Player: players[i]})
			}
		}
	}

	// 按玩家筹码数排序
	sort.SliceStable(winners, func(i, j int) bool {
		return winners[i].Player.Chips < winners[j].Player.Chips
	})
	return winners
}

// GetBestHand 获取玩家的最佳手牌组合。
func GetBestHand(p *Player, t *Table) *Hand {
	// 收集所有可用的牌：2张手牌 + 5张公共牌
	cards := []Card{
		*p.HoleCards[0],
		*p.HoleCards[1],
		*t.Flop[0],
		*t.Flop[1],
		*t.Flop[2],
		*t.Turn,
		*t.River,
	}

	// 找出所有可能的5张牌组合
	// endIndex 是排除法计算得出：总牌数(7) - 手牌数(5) + 1 = 3
	cardCombos := FindCardCombinations(0, 3, cards)

	var cardHand [5]Card
	var bestHand *Hand
	for _, cs := range cardCombos {
		copy(cardHand[:], cs)
		currentHand := CheckHand(cardHand)
		if bestHand == nil {
			// 如果是第一个组合，设为默认最佳手牌
			bestHand = currentHand
		} else {
			result := CompareHand(currentHand, bestHand)
			if result == GreaterThan {
				// 如果当前组合更好，更新最佳手牌
				bestHand = currentHand
			}
		}
	}
	return bestHand
}

// FindCardCombinations 找出所有可能的牌组合。
//
// 在德州扑克中，玩家可以使用 7 张牌（2 张手牌和 5 张公共牌）组成最佳的 5 张牌组合
// 牌的顺序不重要，这将产生21种组合 [即 7! / (5! * 2!)]
//
// 这是一个递归函数，采用暴力方法计算所有可能的组合
// 由于组合数量较小，这种方法是可行的
//
// 初始调用时 startIndex 应为 0
// endIndex 的计算方法是：总牌数 - 手牌数 + 1（因为 endIndex 是排除的）
//
// 示例：7（总牌数）- 5（手牌数）+ 1 = 3
//
// 算法步骤：
// - 遍历每个可能的起始牌
// - 对于每张起始牌，将剩余牌视为子问题继续查找组合
// - 收到子组合后，将起始牌添加到每个子组合中
// - 将所有子组合添加到主组合列表中
func FindCardCombinations(startIndex int, endIndex int, cs []Card) [][]Card {
	// 存储所有可能的牌组合
	combos := make([][]Card, 0)
	numCards := len(cs)

	// 需要 loopCount 因为 i 从 startIndex 开始
	// 这会影响 newStartIndex 的计算，它需要从 0 开始
	loopCount := 0
	for i := startIndex; i < endIndex; i++ {
		// 设置新的起始和结束边界，用于处理剩余的牌子集
		newStartIndex := startIndex + loopCount + 1
		newEndIndex := endIndex + 1

		// 当达到最后一张牌时的基本情况
		subCombos := [][]Card{
			{cs[i]},
		}

		// 如果还没到最后一张牌，递归调用继续查找剩余子集的组合
		if newEndIndex <= numCards {
			subCombos = FindCardCombinations(newStartIndex, newEndIndex, cs)
			// 将当前牌添加到每个子组合中
			for j := range subCombos {
				subCombos[j] = append(subCombos[j], cs[i])
			}
		}

		// 将所有子组合合并到主组合列表中
		combos = append(combos, subCombos...)
		loopCount++
	}

	return combos
}

// CompareHand 比较两副手牌的大小。
func CompareHand(a *Hand, b *Hand) Comparison {
	if a.Rank > b.Rank {
		return GreaterThan
	}

	if a.Rank < b.Rank {
		return LessThan
	}

	// 当两副手牌牌型相同时，使用平局判定规则
	// 由于不同牌型的平局判定规则不同，需要分别处理
	for i := range a.TieBreakers {
		if a.TieBreakers[i] < b.TieBreakers[i] {
			return LessThan
		}
		if a.TieBreakers[i] > b.TieBreakers[i] {
			return GreaterThan
		}
	}

	return EqualTo
}

// CheckHand 检查玩家的手牌组合，判断牌型
func CheckHand(cs [5]Card) *Hand {
	// 按照从大到小的顺序检查可能的牌型
	possibleHands := []IsHand{
		IsRoyalFlush,    // 皇家同花顺
		IsStraightFlush, // 同花顺
		IsFourOfAKind,   // 四条
		IsFullHouse,     // 葫芦
		IsFlush,         // 同花
		IsStraight,      // 顺子
		IsThreeOfAKind,  // 三条
		IsTwoPair,       // 两对
		IsOnePair,       // 一对
		IsHighCard,      // 高牌
	}

	// 从高到低检查每种牌型，返回第一个匹配的结果
	for _, isHand := range possibleHands {
		if hand := isHand(cs); hand != nil {
			return hand
		}
	}
	return nil
}

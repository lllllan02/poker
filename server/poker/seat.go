package poker

import "container/ring"

// Seat 表示扑克牌桌上的一个座位
//
// 实现说明：
// - Seat 内部使用 ring.Ring 实现循环链表结构，便于遍历和定位
// - 座位和节点之间存在双向链接，保证了数据的一致性
//   - 即 Seat.node 包含一个 Ring，其 Value 指向 Seat
//
// - 每个座位可以为空，也可以坐着一个玩家
type Seat struct {
	node   *ring.Ring // 循环链表节点，用于实现座位的循环访问
	Player *Player    // 座位上的玩家，如果座位为空则为 nil
}

// Next 获取桌上的下一个座位
func (s *Seat) Next() *Seat {
	return s.node.Next().Value.(*Seat)
}

// Len 获取桌上座位的总数
func (s *Seat) Len() int {
	return s.node.Len()
}

// NewSeat 创建一个包含 n 个座位的新座位环
func NewSeat(n int) *Seat {
	node := ring.New(n)
	for i := 0; i < node.Len(); i++ {
		node.Value = &Seat{node: node}
		node = node.Next()
	}
	return node.Value.(*Seat)
}

// GetPlayerById 根据玩家 id 查找玩家
func (s *Seat) GetPlayerById(playerId string) *Player {
	for i := 0; i < s.Len(); i++ {
		if s.Player != nil && s.Player.Id == playerId {
			return s.Player
		}
		s = s.Next()
	}
	return nil
}

// GetActivePlayers 获取活跃的玩家
func (s *Seat) GetActivePlayers() []*Player {
	activePlayers := make([]*Player, 0)
	for i := 0; i < s.Len(); i++ {
		if s.Player != nil && s.Player.Status == PlayerActive {
			activePlayers = append(activePlayers, s.Player)
		}
		s = s.Next()
	}
	return activePlayers
}

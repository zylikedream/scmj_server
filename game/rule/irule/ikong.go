package irule

import "zinx-mj/game/card/playercard"

// 杠接口
type IKong interface {
	ExposedKong(pc *playercard.PlayerCard, card int) error   // 明杠
	ConcealedKong(pc *playercard.PlayerCard, card int) error // 暗杠
	Kong(pc *playercard.PlayerCard, card int) error          // 杠牌
}

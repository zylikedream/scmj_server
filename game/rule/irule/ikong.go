package irule

import "zinx-mj/game/card/playercard"

// 杠接口
type IKong interface {
	ExposedKong(pc *playercard.PlayerCard, card int) (*playercard.PlayerCard, error)   // 明杠
	ConcealedKong(pc *playercard.PlayerCard, card int) (*playercard.PlayerCard, error) // 暗杠
	Kong(pc *playercard.PlayerCard, card int) (*playercard.PlayerCard, error)          // 杠牌
}

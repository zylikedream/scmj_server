package irule

import handcard "zinx-mj/game/card/handcard"

// 杠接口
type IKong interface {
	ExposedKong(pc *handcard.HandCard, card int) error   // 明杠
	ConcealedKong(pc *handcard.HandCard, card int) error // 暗杠
	RainKong(pc *handcard.HandCard, card int) error      // 杠牌
}

package irule

import handcard "zinx-mj/game/card/handcard"

// 吃牌接口
type IChow interface {
	Chow(cards *handcard.HandCard, card int) error
}

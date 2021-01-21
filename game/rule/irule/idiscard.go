package irule

import handcard "zinx-mj/game/card/handcard"

// 出牌接口
type IDiscard interface {
	Discard(HandCard *handcard.HandCard, card int) error
}

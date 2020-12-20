package irule

import handcard "zinx-mj/game/card/handcard"

// 摸牌接口
type IDraw interface {
	Draw(pc *handcard.HandCard, card int) error
}

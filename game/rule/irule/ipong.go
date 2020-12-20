package irule

import handcard "zinx-mj/game/card/handcard"

// 碰牌接口
type IPong interface {
	Pong(pc *handcard.HandCard, card int) error
	CanPong(pc *handcard.HandCard, card int) bool
}

package irule

import "zinx-mj/game/card/playercard"

// 出牌接口
type IDiscard interface {
	Discard(playerCard *playercard.PlayerCard, card int, dingqueSuit int) (*playercard.PlayerCard, error)
}

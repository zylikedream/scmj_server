package irule

import "zinx-mj/game/card/playercard"

// 吃牌接口
type IChow interface {
	Chow(cards *playercard.PlayerCard, card int) (*playercard.PlayerCard, error)
}

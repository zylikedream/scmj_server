package irule

import "zinx-mj/game/card/playercard"

// 碰牌接口
type IPong interface {
	Pong(pc *playercard.PlayerCard, card int) error
}

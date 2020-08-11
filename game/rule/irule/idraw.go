package irule

import "zinx-mj/game/card/playercard"

// 摸牌接口
type IDraw interface {
	Draw(pc *playercard.PlayerCard, card int) error
}

package ipoint

import "zinx-mj/game/card/playercard"

// 清一色

type ISingleSuit interface {
	IsSingleSuit(pc *playercard.PlayerCard) bool
}

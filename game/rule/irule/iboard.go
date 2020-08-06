package irule

import "zinx-mj/game/card/boardcard"

// 初牌桌接口
type IBoard interface {
	NewBoard() *boardcard.BoardCard
}

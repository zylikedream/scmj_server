package irule

import "zinx-mj/game/card/boardcard"

type GamePlayer interface {
	IsWin() bool
}

type IGameMode interface {
	IsGameEnd(players []GamePlayer, board *boardcard.BoardCard) bool
}

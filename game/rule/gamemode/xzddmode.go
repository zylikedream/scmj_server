package gamemode

import (
	"zinx-mj/game/card/boardcard"
	"zinx-mj/game/rule/irule"
)

type xzddGameMode struct {
}

func NewXzddGameMode() irule.IGameMode {
	return &xzddGameMode{}
}

func (x *xzddGameMode) IsGameEnd(players []irule.GamePlayer, board *boardcard.BoardCard) bool {
	if board.GetLeftDrawCardNum() == 0 {
		return true
	}
	// 一个人打时要等牌打完

	var winned int
	for _, ply := range players {
		if ply.IsWin() {
			winned++
		}
	}
	if len(players) == 1 { // 1个人特殊处理
		return winned == 1
	}
	return winned == len(players)-1

}

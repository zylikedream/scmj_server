package board

import (
	"zinx-mj/game/card/boardcard"
	"zinx-mj/game/rule/irule"
)

type suitBoard struct {
	suits []int
}

func NewSuitBoard(suits ...int) irule.IBoard {
	return &suitBoard{
		suits: suits,
	}
}

/*
 * Descrp: 初始化筒条万麻将棋牌
 * Create: zhangyi 2020-07-03 11:46:45
 */
func (t *suitBoard) NewBoard() *boardcard.BoardCard {
	// 成都麻将需要筒、条 万
	return boardcard.NewBoardCardBySuit(t.suits...)
}

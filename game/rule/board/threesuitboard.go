package board

import (
	"zinx-mj/game/card"
	"zinx-mj/game/card/boardcard"
	"zinx-mj/game/rule/irule"
)

// 筒条万三种花色的牌
type threeSuitBoard struct {
}

func NewThreeSuitBoard() irule.IBoard {
	return &threeSuitBoard{}
}

/*
 * Descrp: 初始化筒条万麻将棋牌
 * Create: zhangyi 2020-07-03 11:46:45
 */
func (t *threeSuitBoard) NewBoard() *boardcard.BoardCard {
	// 成都麻将需要筒、条 万
	return boardcard.NewBoardCardBySuit(card.CARD_SUIT_DOT, card.CARD_SUIT_BAMBOO, card.CARD_SUIT_CHARACTER)
}

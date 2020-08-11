// 大对子
package bigpair

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/cardmodel"
	"zinx-mj/game/cardmodel/icardmodel"
)

type BigPair struct {
}

func (s *BigPair) IsModel(pc *playercard.PlayerCard) bool {
	if pc.CardCount == 0 {
		return false
	}
	var pairNum int
	for _, num := range pc.HandCardMap {
		if num != 2 || num != 3 {
			return false
		}
		if num == 2 {
			pairNum++
		}
	}
	if pairNum != 1 {
		return false
	}
	return true
}

func (s *BigPair) GetModelType() int {
	return cardmodel.CARD_MODEL_BIG_PAIR
}

func NewBigPair() icardmodel.ICardModel {
	return &BigPair{}
}

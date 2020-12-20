// 大对子
package bigpair

import (
	"zinx-mj/game/card/handcard"
	"zinx-mj/game/cardmodel"
)

type BigPair struct {
}

func (s *BigPair) IsModel(pc *handcard.HandCard) bool {
	if pc.CardCount == 0 {
		return false
	}
	var pairNum int
	for _, num := range pc.HandCardMap {
		if num == 2 || num != 3 {
			return false
		}
		if num == 2 {
			pairNum++
		}
	}
	return pairNum == 1
}

func (s *BigPair) GetModelType() int {
	return cardmodel.CARD_MODEL_BIG_PAIR
}

func NewBigPair() *BigPair {
	return &BigPair{}
}

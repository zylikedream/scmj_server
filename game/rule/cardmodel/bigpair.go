// 大对子
package cardmodel

import (
	"zinx-mj/game/rule/irule"
)

type BigPair struct {
}

func (s *BigPair) IsModel(data *irule.CardModel) bool { // 是否是某种牌型
	var pairNum int
	for _, num := range data.HandCard {
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
	return CARD_MODEL_BIG_PAIR
}

func NewBigPair() *BigPair {
	return &BigPair{}
}

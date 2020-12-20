package sevenpair

import (
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/cardmodel"
	"zinx-mj/game/cardmodel/icardmodel"
)

type SevenPair struct {
}

func (s *SevenPair) IsModel(pc *handcard.HandCard) bool {
	// 不能有碰不能有杠
	if len(pc.KongCards) > 0 || len(pc.PongCards) > 0 {
		return false
	}
	for _, num := range pc.HandCardMap {
		if num%2 != 0 { // 非两个或者四个
			return false
		}
	}
	return true
}

func (s *SevenPair) GetModelType() int {
	return cardmodel.CARD_MODEL_SEVEN_PAIR
}

func NewSevenPair() icardmodel.ICardModel {
	return &SevenPair{}
}

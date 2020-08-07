package sevenpair

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/cardmodel"
	"zinx-mj/game/cardmodel/icardmodel"
)

type SevenPair struct {
}

func (s *SevenPair) IsModel(pc *playercard.PlayerCard) bool {
	// 不能有碰不能有杠
	if len(pc.KongCards) > 0 || len(pc.PongCards) > 0 {
		return false
	}
	for _, num := range pc.HandCardMap {
		if num%2 != 0 { // 非两队或者四队
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
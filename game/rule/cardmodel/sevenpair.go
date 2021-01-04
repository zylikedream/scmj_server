package cardmodel

import (
	"zinx-mj/game/rule/irule"
)

type SevenPair struct {
}

func (s *SevenPair) IsModel(data *irule.CardModel) bool { // 是否是某种牌型
	// 不能有碰不能有杠
	if len(data.KongCard) > 0 || len(data.PongCard) > 0 {
		return false
	}
	for _, num := range data.HandCard {
		if num != 2 { // 7个2对
			return false
		}
	}
	return true
}

func (s *SevenPair) GetModelType() int {
	return CARD_MODEL_SEVEN_PAIR
}

func NewSevenPair() *SevenPair {
	return &SevenPair{}
}

package cardmodel

import (
	"zinx-mj/game/rule/irule"
)

type BigSevenPair struct {
}

func (s *BigSevenPair) IsModel(data *irule.CardModel) bool { // 是否是某种牌型
	// 不能有碰不能有杠
	if len(data.KongCard) > 0 || len(data.PongCard) > 0 {
		return false
	}
	var kongHandCard bool
	for _, num := range data.HandCard {
		if num != 2 && num != 4 { // 非两个或者四个
			return false
		}
		if num == 4 {
			kongHandCard = true
		}
	}
	return kongHandCard // 必须要有4个才是龙七对
}

func (s *BigSevenPair) GetModelType() int {
	return CARD_MODEL_BIG_SEVEN_PAIR
}

func NewBigSevenPair() *BigSevenPair {
	return &BigSevenPair{}
}

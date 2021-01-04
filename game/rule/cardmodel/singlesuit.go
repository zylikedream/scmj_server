// 清一色
package cardmodel

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
)

type SingleSuit struct {
}

func (s *SingleSuit) IsModel(data *irule.CardModel) bool { // 是否是某种牌型
	suit := gamedefine.CARD_SUIT_EMPTY
	for c := range data.HandCard {
		if suit == gamedefine.CARD_SUIT_EMPTY {
			suit = gamedefine.GetCardSuit(c)
		} else if suit != gamedefine.GetCardSuit(c) {
			return false
		}
	}
	for c := range data.KongCard {
		if suit != gamedefine.GetCardSuit(c) {
			return false
		}
	}
	for c := range data.PongCard {
		if suit != gamedefine.GetCardSuit(c) {
			return false
		}
	}
	return true
}

func (s *SingleSuit) GetModelType() int {
	return CARD_MODEL_SINGLE_SUIT
}

func NewSingleSuit() *SingleSuit {
	return &SingleSuit{}
}

// 清一色
package singlesuit

import (
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/cardmodel"
	"zinx-mj/game/cardmodel/icardmodel"
	"zinx-mj/game/gamedefine"
)

type SingleSuit struct {
}

func (s *SingleSuit) IsModel(pc *handcard.HandCard) bool {
	if pc.GetCardTotalCount() == 0 {
		return false
	}
	suit := gamedefine.CARD_SUIT_EMPTY
	for c := range pc.CardMap {
		if suit == gamedefine.CARD_SUIT_EMPTY {
			suit = gamedefine.GetCardSuit(c)
		} else if suit != gamedefine.GetCardSuit(c) {
			return false
		}
	}
	for c := range pc.KongCards {
		if suit != gamedefine.GetCardSuit(c) {
			return false
		}
	}
	for c := range pc.PongCards {
		if suit != gamedefine.GetCardSuit(c) {
			return false
		}
	}
	return true
}

func (s *SingleSuit) GetModelType() int {
	return cardmodel.CARD_MODEL_SINGLE_SUIT
}

func NewSingleSuit() icardmodel.ICardModel {
	return &SingleSuit{}
}

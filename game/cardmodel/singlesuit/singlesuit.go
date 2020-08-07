// 清一色
package singlesuit

import (
	"zinx-mj/game/card"
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/cardmodel"
	"zinx-mj/game/cardmodel/icardmodel"
)

type SingleSuit struct {
}

func (s *SingleSuit) IsModel(pc *playercard.PlayerCard) bool {
	if pc.CardCount == 0 {
		return false
	}
	suit := card.CARD_SUIT_EMPTY
	for c, _ := range pc.HandCardMap {
		if suit == card.CARD_SUIT_EMPTY {
			suit = card.GetCardSuit(c)
		} else if suit != card.GetCardSuit(c) {
			return false
		}
	}
	for c, _ := range pc.KongCards {
		if suit != card.GetCardSuit(c) {
			return false
		}
	}
	for c, _ := range pc.PongCards {
		if suit != card.GetCardSuit(c) {
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

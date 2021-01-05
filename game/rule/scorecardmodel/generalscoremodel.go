package scorecardmodel

import (
	"zinx-mj/game/rule/cardmodel"
	"zinx-mj/game/rule/irule"
)

type generalScoreCardModel struct {
	cardModels []irule.ICardModel
}

func NewGeneralScoreCardModel() *generalScoreCardModel {
	return &generalScoreCardModel{
		[]irule.ICardModel{
			cardmodel.NewBigPair(),
			cardmodel.NewSevenPair(),
			cardmodel.NewBigSevenPair(),
			cardmodel.NewSingleSuit(),
			cardmodel.NewBambooSecond(),
		},
	}
}

func (g *generalScoreCardModel) ScoreCardModel(scard *irule.CardModel) []int {
	var models []int
	for _, cardmodel := range g.cardModels {
		if cardmodel.IsModel(scard) {
			models = append(models, cardmodel.GetModelType())
		}
	}
	return models
}

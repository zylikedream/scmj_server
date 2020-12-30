package icardmodel

import handcard "zinx-mj/game/card/handcard"

type ICardModel interface {
	IsModel(pc *handcard.HandCard) bool // 是否是某种牌型
	GetModelType() int
}

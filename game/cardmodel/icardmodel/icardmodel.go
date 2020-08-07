package icardmodel

import "zinx-mj/game/card/playercard"

type ICardModel interface {
	IsModel(pc *playercard.PlayerCard) bool // 是否是某种牌型
	GetModelType() int
}

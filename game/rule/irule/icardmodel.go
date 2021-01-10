package irule

import "zinx-mj/game/card/handcard"

type CardModel struct {
	HandCard map[int]int
	PongCard []int
	KongCard []handcard.KongInfo
	WinCard  int
	WiRule   IWin
}

type ICardModel interface {
	IsModel(data *CardModel) bool // 是否是某种牌型
	GetModelType() int
}

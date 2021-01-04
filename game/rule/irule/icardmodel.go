package irule

type CardModel struct {
	HandCard map[int]int
	PongCard map[int]struct{}
	KongCard map[int]struct{}
	WinCard  int
	WiRule   IWin
}

type ICardModel interface {
	IsModel(data *CardModel) bool // 是否是某种牌型
	GetModelType() int
}

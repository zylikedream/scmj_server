package irule

// 计算牌型
type IScoreCardModel interface {
	ScoreCardModel(scard *CardModel) []int
}

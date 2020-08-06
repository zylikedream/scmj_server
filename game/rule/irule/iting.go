package irule

type ITing interface {
	// 得到可听的牌
	GetTingCard(cards []int, winRule IWin) []int
	// 是否可以听牌
	CanTing(cards []int, winRule IWin) bool
}

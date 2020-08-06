package irule

// 胡牌接口
type IWin interface {
	CanWin(cards []int) bool
}

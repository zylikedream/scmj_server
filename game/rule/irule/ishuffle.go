package irule

// 洗牌接口
type IShuffle interface {
	Shuffle(cards []int) []int
}

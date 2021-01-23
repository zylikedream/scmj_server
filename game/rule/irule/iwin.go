package irule

const (
	DATA_KEY_DING_QUE = "DING_QUE"
)

// 胡牌接口
type IWin interface {
	CanWin(cards []int) bool
	SetData(key string, value interface{})
}

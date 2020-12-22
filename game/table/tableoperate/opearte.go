package tableoperate

const (
	OPERATE_WIN = iota + 1
	OPERATE_KONG
	OPERATE_PONG
	OPERATE_CHOW
	OPERATE_PASS
	OPERATE_DISCARD
)

type OperateCommand struct {
	OpType int
	OpData OperateData
}

type OperateData struct {
	Card int
}

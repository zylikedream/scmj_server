package tableoperate

const (
	OPERATE_WIN = iota
	OPERATE_KONG
	OPERATE_PONG
	OPERATE_CHOW
	OPERATE_PASS
	OPERATE_DISCARD
)

type PlayerOperate struct {
	Pid    uint64
	OpType int
	OpData OperateData
}

type OperateData struct {
	Card int
}

package tableoperate

const (
	OPERATE_EMPTY = iota // 空操作
	OPERATE_WIN
	OPERATE_KONG_CONCEALED // 暗杠
	OPERATE_KONG_EXPOSED   // 明杠
	OPERATE_KONG_RAIN      // 下雨
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

type PlyOperate struct {
	Pid    uint64
	OpType int
}

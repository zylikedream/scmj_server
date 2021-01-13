package tableoperate

// 按照优先级排序
const (
	OPERATE_EMPTY = iota // 空操作
	OPERATE_WIN
	OPERATE_KONG_CONCEALED // 暗杠
	OPERATE_KONG_EXPOSED   // 明杠
	OPERATE_KONG_RAIN      // 下雨
	OPERATE_PONG
	OPERATE_CHOW
	OPERATE_DISCARD // 出牌
	OPERATE_PASS
)

type OperateCommand struct {
	OpType int
	OpData OperateData
}

type OperateData struct {
	Card int
}

func NewOperateWin(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_WIN,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperateKongConcealed(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_KONG_CONCEALED,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperateKongExposed(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_KONG_EXPOSED,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperateKongRain(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_KONG_RAIN,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperatePong(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_PONG,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperateChow(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_CHOW,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperateDiscard(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_DISCARD,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperatePass() OperateCommand {
	return OperateCommand{
		OpType: OPERATE_PASS,
		OpData: OperateData{
			Card: 0,
		},
	}
}

func NewOperateEmpty() OperateCommand {
	return OperateCommand{
		OpType: OPERATE_EMPTY,
		OpData: OperateData{
			Card: 0,
		},
	}
}

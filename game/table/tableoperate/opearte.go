package tableoperate

// 按照优先级排序
const (
	OPERATE_EMPTY          = iota // 空操作
	OPERATE_DRAW_WIN              // 自摸
	OPERATE_DISCARD_WIN           // 放炮
	OPERATE_KONG_CONCEALED        // 暗杠
	OPERATE_KONG_EXPOSED          // 明杠
	OPERATE_KONG_RAIN             // 下雨
	OPERATE_PONG
	OPERATE_CHOW
	OPERATE_DISCARD  // 出牌
	OPERATE_DING_QUE // 定缺
	OPERATE_DRAW     // 定缺
	OPERATE_PASS
)

type OperateCommand struct {
	OpType int
	OpData OperateData
}

type OperateData struct {
	Card int
}

func NewOperateDrawWin(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_DRAW_WIN,
		OpData: OperateData{
			Card: card,
		},
	}
}

func NewOperateDiscardWin(card int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_DISCARD_WIN,
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

func NewOperateDiscard() OperateCommand {
	return OperateCommand{
		OpType: OPERATE_DISCARD,
		OpData: OperateData{},
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

func NewOperateDingQue() OperateCommand {
	return OperateCommand{
		OpType: OPERATE_DING_QUE,
		OpData: OperateData{},
	}
}

func NewOperateDraw(c int) OperateCommand {
	return OperateCommand{
		OpType: OPERATE_DRAW,
		OpData: OperateData{
			Card: c,
		},
	}
}

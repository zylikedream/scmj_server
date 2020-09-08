package operate

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/player"
)

type OpDiscard struct {
	pid  player.PID
	card int
}

func (s *OpDiscard) GetOperateType() int {
	return gamedefine.OPERATE_DISCARD
}

func NewOperateDiscard() irule.IOperate {
	return &OpDiscard{}
}

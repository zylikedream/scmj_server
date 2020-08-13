package operate

import "zinx-mj/game/rule/irule"

type ScmjOperate struct {
}

func (s *ScmjOperate) GetOperateType() int {
	return 0
}

func NewScmjOperate() irule.IOperate {
	return &ScmjOperate{}
}

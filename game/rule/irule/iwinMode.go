package irule

import "zinx-mj/game/table/tableoperate"

// 计算规则
type IWinMode interface {
	GetWinRule(winPid uint64, turnPid uint64, turnOps []tableoperate.OperateCommand, discards []int) int
}

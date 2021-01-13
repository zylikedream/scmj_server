package irule

import "zinx-mj/game/table/tableoperate"

// 计算规则
type IWinModeModel interface {
	GetWinMode(winPid uint64, turnPid uint64, turnOps []tableoperate.OperateCommand, discards []int) int
}

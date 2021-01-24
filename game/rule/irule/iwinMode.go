package irule

import "zinx-mj/game/table/tableoperate"

// 计算规则
type WinModeInfo struct {
	WinPid   uint64
	DrawWin  bool                          // 自摸
	TurnOps  []tableoperate.OperateCommand // 当前玩家做过的操作
	TurnDraw []int                         // 当前玩家摸过的牌
	Dealer   uint64
	Discards []int
}
type IWinModeModel interface {
	GetWinMode(info WinModeInfo) int
}

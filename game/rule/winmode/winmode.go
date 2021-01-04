package winmode

import "zinx-mj/game/table/tableoperate"

const (
	WIN_MODE_DRAW      = iota // 普通自摸
	WIN_MODE_GOD              // 天胡
	WIN_MODE_KONG_DRAW        // 杠上花

	WIN_MODE_DEVIL        // 地胡
	WIN_MODE_RUB_KONG     // 抢杠
	WIN_MODE_KONG_DISCARD // 杠上炮
	WIN_MODE_DISCARD      // 点炮
)

type WinMode struct {
}

func NewWinMode() *WinMode {
	return &WinMode{}
}

func (w *WinMode) GetWinRule(winPid uint64, turnPid uint64, turnOps []tableoperate.OperateCommand, discards []int) int {
	if winPid == turnPid { // 自摸
		if len(discards) == 0 {
			return WIN_MODE_GOD
		} else {
			lastOp := turnOps[len(turnOps)-1].OpType
			if lastOp == tableoperate.OPERATE_KONG_CONCEALED || lastOp == tableoperate.OPERATE_KONG_EXPOSED || lastOp == tableoperate.OPERATE_KONG_RAIN {
				return WIN_MODE_KONG_DRAW
			}
			return WIN_MODE_DRAW
		}
	} else { // 点炮
		if len(discards) == 1 {
			return WIN_MODE_DEVIL
		} else {
			lastOp := turnOps[len(turnOps)-1].OpType
			if lastOp == tableoperate.OPERATE_KONG_EXPOSED {
				return WIN_MODE_RUB_KONG
			}
			if len(turnOps) >= 2 {
				// 杠上炮的话，当前玩家上上个是杠的动作
				prevOp := turnOps[len(turnOps)-2].OpType
				if prevOp == tableoperate.OPERATE_KONG_EXPOSED || prevOp == tableoperate.OPERATE_KONG_CONCEALED || prevOp == tableoperate.OPERATE_KONG_RAIN {
					return WIN_MODE_KONG_DISCARD
				}
			}
			return WIN_MODE_DISCARD
		}
	}
}

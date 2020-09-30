package sccardtable

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/game/table/tablestate"
	"zinx-mj/player"
)

const (
	TABLE_STATE_DRAW              = "state_draw"
	TABLE_STATE_WAIT_OPERATE      = "state_wait_opearte"
	TABLE_STATE_WAIT_OPERATE_WIN  = "state_wait_opearte_win"
	TABLE_STATE_WAIT_OPERATE_KONG = "state_wait_opearte_kong"
	TABLE_STATE_WAIT_OPERATE_PONG = "state_wait_opearte_pong"
)

type StateDraw struct {
	tablestate.StateBase
}

func NewStateDraw() *StateDraw {
	return &StateDraw{}
}

func (sd *StateDraw) OnEnter(args ...interface{}) error {
	sct := args[0].(*ScCardTable)
	return sct.drawCard(sct.curPlayerIndex)
}

type StateWaitOperate struct {
	tablestate.StateBase
	curOpIndex int
	curPlys    []*tableplayer.TablePlayer
	ops        []int
}

func NewStateWaitOperate() *StateWaitOperate {
	return &StateWaitOperate{
		ops: []int{gamedefine.OPERATE_WIN, gamedefine.OPERATE_KONG, gamedefine.OPERATE_PONG},
	}
}

func (swo *StateWaitOperate) OnEnter(args ...interface{}) error {
	sct := args[0].(*ScCardTable)
	pid := args[1].(player.PID)
	card := args[2].(int)
	swo.curOpIndex = 0
	swo.curPlys = make([]*tableplayer.TablePlayer, 0)

	for _, ply := range sct.players {
		if ply.IsOperateValid(swo.ops[swo.curOpIndex]) {
			swo.curPlys = append(swo.curPlys, ply)
		}
	}
}

func (s *ScCardTable) updateWaitOperateSub(op int, nextState string) error {
	for _, ply := range s.players {
		if ply.IsOperateValid(op) {
			return nil
		}
	}
	return s.stateMachine.Next(nextState)
}

func (s *ScCardTable) enterWaitOperateWin() error {
	return s.enterWaitOperateSub(gamedefine.OPERATE_WIN)
}

func (s *ScCardTable) updateWaitOperateWin() error {
	return s.updateWaitOperateSub(gamedefine.OPERATE_WIN, TABLE_STATE_WAIT_OPERATE_KONG)
}

func (s *ScCardTable) enterWaitOperateKong() error {
	return s.enterWaitOperateSub(gamedefine.OPERATE_KONG)
}

func (s *ScCardTable) updateWaitOperateKong() error {
	return s.updateWaitOperateSub(gamedefine.OPERATE_KONG, TABLE_STATE_WAIT_OPERATE_PONG)
}

func (s *ScCardTable) enterWaitOperatePong() error {
	return s.enterWaitOperateSub(gamedefine.OPERATE_PONG)
}

func (s *ScCardTable) updateWaitOperatePong() error {
	return s.updateWaitOperateSub(gamedefine.OPERATE_PONG, TABLE_STATE_DRAW)
}

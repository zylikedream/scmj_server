package sccardtable

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/player"
)

const (
	TABLE_STATE_DRAW              = "state_draw"
	TABLE_STATE_WAIT_OPERATE      = "state_wait_opearte"
	TABLE_STATE_WAIT_OPERATE_WIN  = "state_wait_opearte_win"
	TABLE_STATE_WAIT_OPERATE_KONG = "state_wait_opearte_kong"
	TABLE_STATE_WAIT_OPERATE_PONG = "state_wait_opearte_pong"
)

func (s *ScCardTable) enterDrawCard() error {
	return s.drawCard(s.curPlayerIndex)
}

func (s *ScCardTable) enterWaitOperate(pid player.PID, card int) error {
	for _, ply := range s.players {
		if ply.Pid == pid {
			continue
		}
		s.FillPlyOperateWithCard(ply, pid, card)
	}
	return nil
}

func (s *ScCardTable) updateWaitOperate() error {
	return s.stateMachine.Next(TABLE_STATE_WAIT_OPERATE_WIN)
}

func (s *ScCardTable) enterWaitOperateSub(op int) error {
	for _, ply := range s.players {
		if ply.IsOperateValid(op) {
			s.NotifyPlyOperate(ply)
		}
	}
	return nil
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

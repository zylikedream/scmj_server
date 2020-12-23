package tablestate

import (
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/util"

	"github.com/pkg/errors"
)

// 明杠有可能被抢杠，所以单独作为一种状态
type StateKongConcealed struct {
	StateBase
	table ITableForState
	pids  []uint64
}

func NewStateKongConcealed(table ITableForState) *StateDiscard {
	return &StateDiscard{
		table: table,
	}
}

func (s *StateKongConcealed) OnEnter() error {
	plys := s.table.GetPlayers()
	for _, ply := range plys {
		ops := ply.GetOperates()
		if len(ops) > 0 {
			s.pids = append(s.pids, ply.Pid)
		}
	}
	return nil
}

func (s *StateKongConcealed) OnUpdate() (IState, error) {
	if len(s.pids) > 0 {
		return nil, nil
	}

	nextState := TABLE_STATE_DRAW
	nextPly := s.table.GetNextTurnPlayer()
	for _, pid := range s.pids {
		if pid == nextPly.Pid {
			// 抢杠只能胡
			nextState = TABLE_STATE_WIN
			break
		}
	}
	return s.stateMachine.GetState(nextState), nil
}

func (s *StateKongConcealed) OnExit() error {
	s.table.UpdateTurnSeat()
	return nil
}

func (s *StateKongConcealed) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	pidIndex := -1
	for i, vpid := range s.pids {
		if vpid == pid {
			pidIndex = i
			break
		}
	}
	if pidIndex == -1 {
		return errors.Errorf("find pid failed, pid=%d, validpids=%v", pid, s.pids)
	}
	util.RemoveElemWithoutOrder(pidIndex, &s.pids)
	return nil
}

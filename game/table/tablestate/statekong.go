package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/pkg/errors"
)

type StateKong struct {
	StateBase
	table ITableForState
	op    int
}

func NewStateKong(table ITableForState) *StateDiscard {
	return &StateDiscard{
		table: table,
	}
}

func (s *StateKong) OnEnter() error {
	return nil
}

func (s *StateKong) OnUpdate() (IState, error) {
	if s.op > 0 {
		return nil, nil
	}

	nextState := getStateByOperate(s.op)
	return s.stateMachine.GetState(nextState), nil
}

func (s *StateKong) OnExit() error {
	s.table.UpdateTurnSeat()
	return nil
}

func (s *StateKong) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	turnPly := s.table.GetTurnPlayer()
	if turnPly.Pid != pid {
		return errors.Errorf("player error, want:%d, get:%d", turnPly.Pid, pid)
	}
	if !turnPly.IsOperateValid(data.OpType) {
		return errors.Errorf("operate unvalid, op:%d validops:%v", data.OpType, turnPly.GetOperates())
	}
	s.op = data.OpType
	return nil
}

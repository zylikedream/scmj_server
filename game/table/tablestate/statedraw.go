package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/pkg/errors"
)

type StateDraw struct {
	StateBase
	table ITableForState
	op    int
}

func NewStateDraw(table ITableForState) *StateDraw {
	return &StateDraw{
		table: table,
		op:    tableoperate.OPERATE_EMPTY,
	}
}

func (s *StateDraw) OnEnter() error {
	if err := s.table.DrawCard(); err != nil {
		return err
	}
	return nil
}

func (s *StateDraw) OnUpdate() (IState, error) {
	if s.op == tableoperate.OPERATE_EMPTY {
		return nil, nil
	}

	nextState := getStateByOperate(s.op)
	return s.stateMachine.GetState(nextState), nil
}

func (s *StateDraw) OnExit() error {
	s.table.UpdateTurnSeat()
	s.Reset()
	return nil
}

func (s *StateDraw) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	turnPly := s.table.GetTurnPlayer()
	if turnPly.Pid != pid {
		return errors.Errorf("player error, want:%d, get:%d", turnPly.Pid, pid)
	}
	if data.OpType == tableoperate.OPERATE_PASS { // 跳过pass操作
		return nil
	}
	s.op = data.OpType
	return nil
}

func (s *StateDraw) Reset() {
	s.op = tableoperate.OPERATE_EMPTY
}

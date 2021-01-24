package tablestate

import (
	"zinx-mj/game/table/tableoperate"
)

type StatePong struct {
	StateBase
	table ITableForState
	op    int
}

func NewStatePong(table ITableForState) *StatePong {
	return &StatePong{
		table: table,
		op:    tableoperate.OPERATE_EMPTY,
	}
}

func (s *StatePong) OnEnter() error {
	return nil
}

func (s *StatePong) OnUpdate() (IState, error) {
	if s.op == tableoperate.OPERATE_EMPTY {
		return nil, nil
	}
	nextState := getStateByOperate(s.op)
	return s.stateMachine.GetState(nextState), nil
}

func (s *StatePong) OnExit() error {
	s.table.UpdateTurnSeat()
	s.op = tableoperate.OPERATE_EMPTY
	return nil
}

func (s *StatePong) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	s.op = data.OpType
	return nil
}

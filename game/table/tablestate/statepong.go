package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/pkg/errors"
)

type StatePong struct {
	StateBase
	table ITableForState
}

func NewStatePong(table ITableForState) *StateDiscard {
	return &StateDiscard{
		table: table,
	}
}

func (s *StatePong) OnEnter() error {
	return nil
}

func (s *StatePong) OnUpdate() (IState, error) {
	return s.stateMachine.GetState(TABLE_STATE_DISCARD), nil
}

func (s *StatePong) OnExit() error {
	s.table.UpdateTurnSeat()
	return nil
}

func (s *StatePong) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	return errors.Errorf("%s state not allow any operate", s.name)
}

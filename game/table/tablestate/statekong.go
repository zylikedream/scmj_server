package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/pkg/errors"
)

type StateKong struct {
	StateBase
	table ITableForState
}

func NewStateKong(table ITableForState) *StateKong {
	return &StateKong{
		table: table,
	}
}

func (s *StateKong) OnEnter() error {
	return nil
}

func (s *StateKong) OnUpdate() (IState, error) {
	return s.stateMachine.GetState(TABLE_STATE_DRAW), nil
}

func (s *StateKong) OnExit() error {
	s.table.UpdateTurnSeat()
	return nil
}

func (s *StateKong) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	return errors.Errorf("%s state not allow any operate", s.name)
}

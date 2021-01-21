package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/pkg/errors"
)

type StateInit struct {
	StateBase
	table ITableForState
}

func NewStateInit(table ITableForState) *StateInit {
	return &StateInit{
		table: table,
	}
}

func (s *StateInit) OnEnter() error {
	return nil
}

func (s *StateInit) OnUpdate() (IState, error) {
	if !s.table.IsGameStart() {
		return nil, nil
	}
	return s.stateMachine.GetState(TABLE_STATE_DING_QUE), nil
}

func (s *StateInit) OnExit() error {
	return nil
}

func (s *StateInit) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	return errors.Errorf("%s state not allow any operate", s.name)
}

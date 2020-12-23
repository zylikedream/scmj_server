package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/pkg/errors"
)

type StateWin struct {
	StateBase
	table ITableForState
}

func NewStateWin(table ITableForState) *StateDiscard {
	return &StateDiscard{
		table: table,
	}
}

func (s *StateWin) OnEnter() error {
	s.table.SetNextSeat(s.table.GetTurnSeat() + 1) // 轮转到下一位
	return nil
}

func (s *StateWin) OnUpdate() (IState, error) {
	// 下一位肯定是摸牌
	return s.stateMachine.GetState(TABLE_STATE_DRAW), nil
}

func (s *StateWin) OnExit() error {
	s.table.UpdateTurnSeat()
	return nil
}

func (s *StateWin) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	return errors.Errorf("%s state not allow any operate", s.name)
}

package tablestate

import (
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/network/protocol"
)

type StateDingQue struct {
	StateBase
	table   ITableForState
	dingQue map[uint64]int
}

func NewStateDingQue(table ITableForState) *StateDingQue {
	return &StateDingQue{
		table:   table,
		dingQue: make(map[uint64]int),
	}
}

func (s *StateDingQue) OnEnter() error {
	// 通知定缺
	for _, ply := range s.table.GetPlayers() {
		ply.SetOperate([]tableoperate.OperateCommand{tableoperate.NewOperateDingQue()})
	}
	msg := &protocol.ScDingQueStart{}
	for _, suit := range s.table.GetBoardCard().GetSuits() {
		msg.Suit = append(msg.Suit, int32(suit))
	}
	_ = s.table.BroadCast(protocol.PROTOID_SC_DING_QUE_START, msg)
	return nil
}

func (s *StateDingQue) OnUpdate() (IState, error) {
	if len(s.dingQue) < len(s.table.GetPlayers()) {
		return nil, nil
	}
	return s.stateMachine.GetState(TABLE_STATE_DRAW), nil
}

func (s *StateDingQue) OnExit() error {
	msg := &protocol.ScDingQueFinish{}
	for pid, suit := range s.dingQue {
		dqInfo := &protocol.ScDingqueInfo{
			Pid:    pid,
			DqSuit: int32(suit),
		}
		msg.Dingque = append(msg.Dingque, dqInfo)
	}
	_ = s.table.BroadCast(protocol.PROTOID_SC_DING_QUE_FINISH, msg)

	s.dingQue = make(map[uint64]int)
	return nil
}

func (s *StateDingQue) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	s.dingQue[pid] = data.OpData.Card
	return nil
}

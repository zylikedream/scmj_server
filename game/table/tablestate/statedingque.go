package tablestate

import (
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/network/protocol"

	"github.com/pkg/errors"
)

type StateDingQue struct {
	StateBase
	draw    *StateDraw
	table   ITableForState
	dingQue map[uint64]int
}

func NewStateDingQue(table ITableForState) *StateDingQue {
	return &StateDingQue{
		table:   table,
		dingQue: make(map[uint64]int),
		draw:    NewStateDraw(table),
	}
}

func (s *StateDingQue) SetStateMachine(sm IStateMachine) {
	s.stateMachine = sm
	s.draw.stateMachine = sm
}

func (s *StateDingQue) IsDingQueFinish() bool {
	return len(s.dingQue) >= len(s.table.GetPlayers())
}

func (s *StateDingQue) OnEnter() error {
	if err := s.draw.OnEnter(); err != nil {
		return err
	}
	for _, ply := range s.table.GetPlayers() {
		ply.SetOperate([]tableoperate.OperateCommand{tableoperate.NewOperateDingQue()})
	}
	// 通知定缺
	msg := &protocol.ScDingQueStart{}
	for _, suit := range s.table.GetBoardCard().GetSuits() {
		msg.Suit = append(msg.Suit, int32(suit))
	}
	_ = s.table.BroadCast(protocol.PROTOID_SC_DING_QUE_START, msg)
	return nil
}

func (s *StateDingQue) OnUpdate() (IState, error) {
	return s.draw.OnUpdate()
}

func (s *StateDingQue) OnExit() error {
	s.dingQue = make(map[uint64]int)
	return s.draw.OnExit()
}

func (s *StateDingQue) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	if s.IsDingQueFinish() {
		return s.draw.OnPlyOperate(pid, data)
	}
	if data.OpType != tableoperate.OPERATE_DING_QUE {
		return errors.Errorf("can't allow operate in dingque, op:%d", data.OpType)
	}
	s.dingQue[pid] = data.OpData.Card
	if !s.IsDingQueFinish() {
		return nil
	}
	// 所有玩家定缺成功

	if err := s.table.DingQueFinish(); err != nil {
		return err
	}
	return nil
}

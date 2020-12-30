package tablestate

import (
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/network/protocol"
	"zinx-mj/util"

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
	// 按照优先级排序
	if err := s.table.DrawCard(); err != nil {
		return err
	}
	if err := s.distributeOperate(); err != nil {
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
	if !turnPly.IsOperateValid(data.OpType) {
		return errors.Errorf("operate unvalid, op:%d validops:%v", data.OpType, turnPly.GetOperates())
	}
	s.op = data.OpType
	return nil
}

func (s *StateDraw) Reset() {
	s.op = tableoperate.OPERATE_EMPTY
}

func (s *StateDraw) distributeOperate() error {
	turnPly := s.table.GetTurnPlayer()
	ops := turnPly.GetOperates()
	if len(ops) == 0 {
		return nil
	}
	opdata := &protocol.ScPlayerOperate{}
	for _, op := range ops {
		opdata.OpType = append(opdata.OpType, int32(op))
	}
	if err := util.SendMsg(turnPly.Pid, protocol.PROTOID_SC_PLAYER_OPERATE, opdata); err != nil {
		return err
	}
	return nil
}

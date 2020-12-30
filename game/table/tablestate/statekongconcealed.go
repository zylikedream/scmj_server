package tablestate

import (
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/network/protocol"
	"zinx-mj/util"

	"github.com/pkg/errors"
)

// 明杠有可能被抢杠，所以单独作为一种状态
type StateKongConcealed struct {
	StateBase
	table ITableForState
	pids  []uint64
}

func NewStateKongConcealed(table ITableForState) *StateKongConcealed {
	return &StateKongConcealed{
		table: table,
	}
}

func (s *StateKongConcealed) Reset() {
	s.pids = s.pids[0:0]
}

func (s *StateKongConcealed) OnEnter() error {
	plys := s.table.GetPlayers()
	for _, ply := range plys {
		ops := ply.GetOperates()
		if len(ops) > 0 {
			s.pids = append(s.pids, ply.Pid)
		}
	}
	if err := s.distributeOperate(); err != nil {
		return err
	}
	return nil
}

func (s *StateKongConcealed) OnUpdate() (IState, error) {
	if len(s.pids) > 0 { // 等待操作
		return nil, nil
	}

	nextState := TABLE_STATE_DRAW
	nextPly := s.table.GetNextTurnPlayer()
	for _, pid := range s.pids {
		if pid == nextPly.Pid {
			// 抢杠只能胡
			nextState = TABLE_STATE_WIN
			break
		}
	}
	return s.stateMachine.GetState(nextState), nil
}

func (s *StateKongConcealed) OnExit() error {
	s.table.UpdateTurnSeat()
	s.Reset()
	return nil
}

func (s *StateKongConcealed) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	pidIndex := -1
	for i, vpid := range s.pids {
		if vpid == pid {
			pidIndex = i
			break
		}
	}
	if pidIndex == -1 {
		return errors.Errorf("find pid failed, pid=%d, validpids=%v", pid, s.pids)
	}
	if data.OpType != tableoperate.OPERATE_WIN {
		return errors.Errorf("kong concealed can only do win operate, pid=%d, data=%v", pid, data)
	}
	util.RemoveElemWithoutOrder(pidIndex, &s.pids)
	return nil
}

func (s *StateKongConcealed) distributeOperate() error {
	opdata := &protocol.ScPlayerOperate{
		OpType: []int32{tableoperate.OPERATE_WIN, tableoperate.OPERATE_PASS},
	}
	for _, pid := range s.pids {
		if err := util.SendMsg(pid, protocol.PROTOID_SC_PLAYER_OPERATE, opdata); err != nil {
			return err
		}
	}
	return nil
}

package tablestate

import (
	"sort"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/util"

	"github.com/pkg/errors"
)

type Action struct {
	Op   int
	Pids []uint64
}

type OpLog struct {
	Op  int
	Pid uint64
}

type StateDiscard struct {
	StateBase
	table ITableForState
	acts  []Action
	oplog []OpLog
}

var opOrder = []int{tableoperate.OPERATE_WIN, tableoperate.OPERATE_KONG_WIND, tableoperate.OPERATE_PONG}

func NewStateDiscard(table ITableForState) *StateDiscard {
	return &StateDiscard{
		table: table,
	}
}

func (s *StateDiscard) Reset() {
	s.acts = s.acts[0:0]
	s.oplog = s.oplog[0:0]
}

func (s *StateDiscard) OnEnter() error {
	s.table.SetNextSeat(s.table.GetTurnSeat() + 1) // 默认轮转到下一位
	// 初始化玩家的操作
	plys := s.table.GetPlayers()
	for _, ply := range plys {
		ops := ply.GetOperates()
		if len(ops) > 0 {
			s.addActions(ply.Pid, ops)
		}
	}
	// 按照优先级排序
	if len(s.acts) > 1 {
		sort.Slice(s.acts, func(i, j int) bool {
			return s.acts[i].Op < s.acts[j].Op
		})
	}

	return nil
}

func (s *StateDiscard) addActions(pid uint64, ops []int) {
	for _, op := range opOrder {
		var find bool
		for _, act := range s.acts {
			if act.Op == op {
				find = true
				act.Pids = append(act.Pids, pid)
				break
			}
		}
		if !find {
			s.acts = append(s.acts, Action{Op: op, Pids: []uint64{pid}})
		}
	}
}

func (s *StateDiscard) OnUpdate() (IState, error) {
	if len(s.acts) > 0 {
		return nil, nil
	}
	nextPly := s.table.GetNextTurnPlayer()
	// 得到下一个玩家的操作, 如果没有默认就是抽牌
	var op int
	for _, log := range s.oplog {
		if log.Pid == nextPly.Pid {
			op = log.Op
		}
	}
	nextState := getStateByOperate(op) // 得到对应状态
	return s.stateMachine.GetState(nextState), nil
}

func (s *StateDiscard) OnExit() error {
	s.table.UpdateTurnSeat()
	s.Reset()
	return nil
}

func (s *StateDiscard) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	if len(s.acts) == 0 {
		return errors.Errorf("can't operate, no acts, op:%v state:%s", data, s.name)
	}
	curAct := s.acts[0] // 当前等待的操作

	// 操作不符且不是跳过
	if curAct.Op != data.OpType && data.OpType != tableoperate.OPERATE_PASS {
		return errors.Errorf("can't operate, wait another act, act:%v op:%v", curAct, data)
	}
	// 查找该玩家是否能操作
	pidIndex := -1
	for i, actPid := range curAct.Pids {
		if actPid == pid {
			pidIndex = i
			break
		}
	}
	if pidIndex == -1 {
		return errors.Errorf("can't operate, act has no pid, act:%v, op:%v", curAct, data)
	}

	// 删除该pids
	util.RemoveElemWithoutOrder(pidIndex, &curAct.Pids)

	s.oplog = append(s.oplog, OpLog{
		Op:  data.OpType,
		Pid: pid,
	})
	if len(curAct.Pids) == 0 { // 一牌不能多用, 直接返回
		return nil
	}

	return nil
}

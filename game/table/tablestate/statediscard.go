package tablestate

import (
	"zinx-mj/game/table/tableoperate"

	"github.com/aceld/zinx/zlog"
)

type OpLog struct {
	Op  int
	Pid uint64
}

type StateDiscard struct {
	StateBase
	table        ITableForState
	oplog        []OpLog
	plyOperate   map[uint64][]int
	operateCount map[int]int
}

func NewStateDiscard(table ITableForState) *StateDiscard {
	return &StateDiscard{
		table:        table,
		plyOperate:   make(map[uint64][]int),
		operateCount: make(map[int]int),
	}
}

func (s *StateDiscard) Reset() {
	s.oplog = s.oplog[0:0]
	s.plyOperate = make(map[uint64][]int)
	s.operateCount = make(map[int]int)
}

func (s *StateDiscard) OnEnter() error {
	s.table.SetNextSeat(s.table.GetTurnSeat() + 1) // 默认轮转到下一位
	// 初始化玩家的操作
	for _, ply := range s.table.GetPlayers() {
		s.addOperates(ply.Pid, ply.GetOperates())
	}
	return nil
}

func (s *StateDiscard) addOperates(pid uint64, ops []int) {
	if len(ops) == 0 {
		return
	}
	s.plyOperate[pid] = ops
	for _, op := range ops {
		s.operateCount[op] += 1
	}
}

func (s *StateDiscard) getOpLog(pid uint64) int {
	op := tableoperate.OPERATE_EMPTY
	for _, log := range s.oplog {
		if log.Pid == pid {
			op = log.Op
			break
		}
	}
	return op
}

func (s *StateDiscard) OnUpdate() (IState, error) {
	if len(s.plyOperate) > 0 {
		return nil, nil
	}
	nextPly := s.table.GetNextTurnPlayer()
	// 得到下一个玩家的操作, 如果没有默认就是抽牌
	op := s.getOpLog(nextPly.Pid)
	nextState := getStateByOperate(op) // 得到对应状态
	return s.stateMachine.GetState(nextState), nil
}

func (s *StateDiscard) OnExit() error {
	s.table.UpdateTurnSeat()
	s.Reset()
	return nil
}

func (s *StateDiscard) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	s.oplog = append(s.oplog, OpLog{ // 保存玩家的操作记录
		Op:  data.OpType,
		Pid: pid,
	})

	if data.OpType == tableoperate.OPERATE_PASS {
		delete(s.plyOperate, pid) // 删除玩家的操作
		return nil
	}
	for opid, ops := range s.plyOperate {
		if opid == pid {
			continue
		}
		if ops[0] < data.OpType { // 如果有更高优先级的操作那么，就需要等待
			zlog.Warnf("need to wait operates, pid:%d, op:%d, waitpid:%d, ops:%v", pid, data.OpType, opid, ops)
			return nil
		}
	}

	delete(s.plyOperate, pid) // 删除玩家的操作
	s.operateCount[data.OpType] -= 1
	if s.operateCount[data.OpType] > 0 { // 还有同等优先级的操作，需要等待其他
		return nil
	}
	// 一牌不能多用, 其他操作不能再使用了, 直接清空所有可以做的操作
	s.plyOperate = make(map[uint64][]int)
	return nil
}

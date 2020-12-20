package tablestate

import (
	"fmt"
	"sort"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/game/table/tableplayer"

	"github.com/pkg/errors"
)

type ITableForState interface {
	GetPlayers() []*tableplayer.TablePlayer
}

const (
	TABLE_STATE_DRAW    = "state_draw"
	TABLE_STATE_DISCARD = "state_discard"
)

type StateMachine struct {
	curState IState
	states   map[string]IState
}

func New(table ITableForState) *StateMachine {
	return &StateMachine{}
}

func (sm *StateMachine) GetCurState() IState {
	return sm.curState
}

func (sm *StateMachine) SetInitState(stateName string) error {
	if state, ok := sm.states[stateName]; !ok {
		return errors.Errorf("no state, name:%s", stateName)
	} else {
		state.OnEnter()
	}
	return nil
}

// Trigger trigger an event
func (sm *StateMachine) Next(nextState IState, value ...interface{}) error {
	// State: exit
	if sm.curState != nil {
		if err := sm.curState.OnExit(); err != nil {
			return fmt.Errorf("onExit state %s err, err=%s", sm.curState.GetName(), err)
		}
	}

	sm.curState = nextState
	if err := nextState.OnEnter(value...); err != nil {
		return fmt.Errorf("onEnter state %s err, err=%s", nextState.GetName(), err)
	}
	return nil
}

func (sm *StateMachine) Update() error {
	nextState, err := sm.curState.OnUpdate()
	if err != nil {
		return err
	}
	return sm.Next(nextState)
}

type IState interface {
	GetName() string
	OnEnter(args ...interface{}) error
	OnUpdate(args ...interface{}) (IState, error)
	OnPlyOperate(data tableoperate.PlayerOperate) error
	OnExit() error
}

type StateBase struct {
	name string
}

func (s *StateBase) GetName() string {
	return s.name
}

func (s *StateBase) OnEnter(args ...interface{}) error {
	return nil
}

func (s *StateBase) OnUpdate(args ...interface{}) (IState, error) {
	return nil, nil
}

func (s *StateBase) OnExit() error {
	return nil
}

func (s *StateBase) OnPlyOperate(data tableoperate.PlayerOperate) error {
	return nil
}

type Action struct {
	Op   int
	Pids []uint64
}

type StateDiscard struct {
	StateBase
	table     ITableForState
	Acts      []Action
	nextState IState
	opOrder   []int
}

func NewStateDiscard(table ITableForState) *StateDiscard {
	opOrder := []int{tableoperate.OPERATE_WIN, tableoperate.OPERATE_KONG, tableoperate.OPERATE_PONG}
	return &StateDiscard{
		table:   table,
		opOrder: opOrder,
		Acts:    make([]Action, len(opOrder)),
	}
}

func (s *StateDiscard) OnEnter(args ...interface{}) error {
	pid := args[0].(uint64)
	c := args[1].(int)
	plys := s.table.GetPlayers()
	for i := range plys {
		ply := plys[i]
		if ply.Pid == pid {
			continue
		}
		ops := ply.GetOperateOnOtherTurn(c)
		ply.AddOperate(ops...)
		s.AddActions(pid, ops)
	}
	// 按照优先级排序
	sort.Slice(s.Acts, func(i, j int) bool {
		return s.Acts[i].Op < s.Acts[j].Op
	})
	return nil
}

func (s *StateDiscard) AddActions(pid uint64, ops []int) {
	for _, op := range s.opOrder {
		var find bool
		for _, act := range s.Acts {
			if act.Op == op {
				find = true
				act.Pids = append(act.Pids, pid)
				break
			}
		}
		if !find {
			s.Acts = append(s.Acts, Action{Op: op, Pids: []uint64{pid}})
		}
	}
}

func (s *StateDiscard) OnUpdate(args ...interface{}) (IState, error) {
	if len(s.Acts) == 0 {
		return s.nextState, nil
	}
	return nil, nil
}

func (s *StateDiscard) OnPlyOperate(data tableoperate.PlayerOperate) error {
	if len(s.Acts) == 0 {
		return errors.Errorf("can't operate, no acts, op:%v state:%s", data, s.name)
	}
	curAct := s.Acts[0] // 当前等待的操作

	// 操作不符且不是跳过
	if curAct.Op != data.OpType && data.OpType != tableoperate.OPERATE_PASS {
		return errors.Errorf("can't operate, wait another act, act:%v op:%v", curAct, data)
	}
	// 查找该玩家是否能操作
	pidIndex := -1
	for i, pid := range curAct.Pids {
		if pid == data.Pid {
			pidIndex = i
			break
		}
	}
	if pidIndex == -1 {
		return errors.Errorf("can't operate, act has no pid, act:%v, op:%v", curAct, data)
	}

	// 删除该pids
	curAct.Pids[pidIndex] = 0
	// 是否还有需要等待的玩家
	for _, pid := range curAct.Pids {
		if pid != 0 {
			return nil
		}
	}

	return nil
}

package tablestate

import (
	"fmt"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/game/table/tableplayer"

	"github.com/aceld/zinx/zlog"

	"github.com/pkg/errors"
)

type ITableForState interface {
	UpdateTurnSeat()
	GetTurnPlayer() *tableplayer.TablePlayer
	DrawCard() error
	GetPlayers() []*tableplayer.TablePlayer
	GetTurnSeat() int
	SetNextSeat(seat int)
	GetNextTurnPlayer() *tableplayer.TablePlayer
	IsGameStart() bool
}

const (
	TABLE_STATE_WIN            = "state_win"
	TABLE_STATE_KONG           = "state_kong"
	TABLE_STATE_KONG_CONCEALED = "state_kong_concealed"
	TABLE_STATE_PONG           = "state_gong"
	TABLE_STATE_DISCARD        = "state_discard"
	TABLE_STATE_DRAW           = "state_draw"
	TABLE_STATE_INIT           = "state_init"
)

type StateMachine struct {
	curState IState
	states   map[string]IState
}

func New(table ITableForState) *StateMachine {
	sm := &StateMachine{
		states: map[string]IState{
			TABLE_STATE_DRAW:           NewStateDraw(table),
			TABLE_STATE_DISCARD:        NewStateDiscard(table),
			TABLE_STATE_WIN:            NewStateWin(table),
			TABLE_STATE_KONG:           NewStateKong(table),
			TABLE_STATE_KONG_CONCEALED: NewStateKongConcealed(table),
			TABLE_STATE_PONG:           NewStatePong(table),
			TABLE_STATE_INIT:           NewStateInit(table),
		},
	}
	for name, state := range sm.states {
		state.SetStateMachine(sm)
		state.SetName(name)
	}
	return sm
}

func (sm *StateMachine) GetCurState() IState {
	return sm.curState
}

func (sm *StateMachine) GetState(name string) IState {
	return sm.states[name]
}

func (sm *StateMachine) SetInitState(stateName string) error {
	if state, ok := sm.states[stateName]; !ok {
		return errors.Errorf("no state, name:%s", stateName)
	} else {
		sm.curState = state
		if err := state.OnEnter(); err != nil {
			return err
		}
	}
	return nil
}

func (sm *StateMachine) Next(nextState IState) error {
	if sm.curState != nil {
		if err := sm.curState.OnExit(); err != nil {
			return fmt.Errorf("onExit state %s err, err=%s", sm.curState.GetName(), err)
		}
	}
	zlog.Infof("change state: %s->%s", sm.curState.GetName(), nextState.GetName())
	sm.curState = nextState
	if err := nextState.OnEnter(); err != nil {
		return fmt.Errorf("onEnter state %s err, err=%s", nextState.GetName(), err)
	}
	return nil
}

func (sm *StateMachine) Update() error {
	if sm.curState == nil {
		return nil
	}
	nextState, err := sm.curState.OnUpdate()
	if err != nil {
		return err
	}
	if nextState == nil {
		return nil
	}
	return sm.Next(nextState)
}

func getStateByOperate(op int) string {
	switch op {
	case tableoperate.OPERATE_WIN:
		return TABLE_STATE_WIN
	case tableoperate.OPERATE_KONG_WIND, tableoperate.OPERATE_KONG_EXPOSED:
		return TABLE_STATE_KONG
	case tableoperate.OPERATE_KONG_CONCEALED:
		return TABLE_STATE_KONG_CONCEALED
	case tableoperate.OPERATE_PONG:
		return TABLE_STATE_PONG
	case tableoperate.OPERATE_DISCARD:
		return TABLE_STATE_DISCARD
	default: // 默认是抽牌
		return TABLE_STATE_DRAW
	}
}

type IStateMachine interface {
	GetState(name string) IState
}

type IState interface {
	GetName() string
	SetName(name string)
	OnEnter() error
	OnUpdate() (IState, error)
	OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error
	OnExit() error
	SetStateMachine(sm IStateMachine)
	Reset()
}

type StateBase struct {
	name         string
	stateMachine IStateMachine
}

func (s *StateBase) GetName() string {
	return s.name
}

func (s *StateBase) SetName(name string) {
	s.name = name
}

func (s *StateBase) OnEnter() error {
	return nil
}

func (s *StateBase) OnUpdate() (IState, error) {
	return nil, nil
}

func (s *StateBase) OnExit() error {
	s.Reset()
	return nil
}

func (s *StateBase) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	return nil
}

func (s *StateBase) SetStateMachine(sm IStateMachine) {
	s.stateMachine = sm
}

func (s *StateBase) Reset() {
}

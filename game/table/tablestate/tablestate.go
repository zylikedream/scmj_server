package tablestate

import (
	"fmt"
)

const (
	TABLE_STATE_DRAW              = "state_draw"
	TABLE_STATE_WAIT_OPERATE      = "state_wait_opearte"
	TABLE_STATE_WAIT_OPERATE_WIN  = "state_wait_opearte_win"
	TABLE_STATE_WAIT_OPERATE_KONG = "state_wait_opearte_kong"
	TABLE_STATE_WAIT_OPERATE_PONG = "state_wait_opearte_pong"
)

type StateMachine struct {
	curState IState
	states   map[string]IState
}

func New() *StateMachine {
	return &StateMachine{
		states: make(map[string]IState),
	}
}

func (sm *StateMachine) GetCurStateName() string {
	if sm.curState == nil {
		return ""
	}
	return sm.curState.GetName()
}

// State define a state
func (sm *StateMachine) AddState(state IState) {
	sm.states[state.GetName()] = state
}

// Trigger trigger an event
func (sm *StateMachine) Next(name string, value ...interface{}) error {
	// State: exit
	if sm.curState != nil {
		if err := sm.curState.OnExit(); err != nil {
			return fmt.Errorf("onExit state %s err, err=%s", sm.curState.GetName(), err)
		}
	}

	if state, ok := sm.states[name]; ok {
		if err := state.OnEnter(value...); err != nil {
			return fmt.Errorf("onEnter state %s err, err=%s", name, err)
		}
	}
	return nil
}

type IState interface {
	GetName() string
	OnEnter(args ...interface{}) error
	OnUpdate(args ...interface{}) error
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

func (s *StateBase) OnUpdate(args ...interface{}) error {
	return nil
}

func (s *StateBase) OnExit() error {
	return nil
}

package tablestate

import (
	"fmt"
	"reflect"

	"github.com/aceld/zinx/zlog"
)

type StateMachine struct {
	curState *State
	states   map[string]*State
}

func New() *StateMachine {
	return &StateMachine{
		states: make(map[string]*State),
	}
}

func (sm *StateMachine) GetCurStateName() string {
	if sm.curState == nil {
		return ""
	}
	return sm.curState.Name
}

func (sm *StateMachine) GetCurState() *State {
	return sm.curState
}

func (sm *StateMachine) InitalState(name string, value ...interface{}) error {
	sm.curState = sm.states[name]
	if sm.curState == nil {
		return fmt.Errorf("cant' find initalstate %s", name)
	}
	return sm.curState.onEnter(value...)
}

// State define a state
func (sm *StateMachine) State(name string) *State {
	state := &State{Name: name}
	sm.states[name] = state
	return state
}

// Trigger trigger an event
func (sm *StateMachine) Next(name string, value ...interface{}) error {
	// State: exit
	if sm.curState != nil {
		if err := sm.curState.onExit(); err != nil {
			return fmt.Errorf("onExit state %s err, err=%s", sm.curState.Name, err)
		}
	}

	if state, ok := sm.states[name]; ok {
		if err := state.onEnter(value...); err != nil {
			return fmt.Errorf("onEnter state %s err, err=%s", name, err)
		}
	}
	return nil
}

type stateFunc func(...interface{}) error
type State struct {
	onEnter  stateFunc
	onUpdate func() error
	onExit   func() error
	Name     string
}

func newState(name string) *State {
	return &State{
		Name: name,
		onEnter: func(value ...interface{}) error {
			return nil
		},
		onExit: func() error {
			return nil
		},
	}
}

func (s *State) Enter(f interface{}) *State {
	var err error
	if s.onEnter, err = s.decorator(f); err != nil {
		zlog.Errorf("set Enter func failed")
		return nil
	}
	return s
}

func (s *State) Exit(f func() error) *State {
	s.onExit = f
	return s
}

func (s *State) Update(f func() error) *State {
	s.onUpdate = f
	return s
}

func (s *State) decorator(f interface{}) (stateFunc, error) {
	ftype := reflect.TypeOf(f)
	if ftype.Kind() != reflect.Func {
		return nil, fmt.Errorf("need a func")
	}
	return func(args ...interface{}) error {
		if ftype.NumIn() != len(args) {
			return fmt.Errorf("on enter failed, params need=%d, get=%d", ftype.NumIn(), len(args))
		}
		if ftype.NumOut() != 1 {
			return fmt.Errorf("func must return an error")
		}
		if !ftype.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return fmt.Errorf("func must return an error")
		}
		params := []reflect.Type{}
		for i := 0; i < ftype.NumIn(); i++ {
			params = append(params, ftype.In(i))
		}
		argValues := make([]reflect.Value, 0, len(params))
		for i := 0; i < len(args); i++ {
			argValues = append(argValues, reflect.ValueOf(args[i]))
		}
		res := reflect.ValueOf(f).Call(argValues)
		if err, ok := (res[0].Interface()).(error); ok {
			return err
		}
		return nil
	}, nil
}

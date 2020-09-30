package tablestate

import (
	"fmt"
	"reflect"

	"github.com/aceld/zinx/zlog"
)

type StateMachine struct {
	curState *State
	states   map[string]*State
	events   map[string]*Event
}

func New() *StateMachine {
	return &StateMachine{
		states: make(map[string]*State),
		events: make(map[string]*Event),
	}
}

func (sm *StateMachine) GetCurStateName() string {
	if sm.curState == nil {
		return ""
	}
	return sm.curState.Name
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

// Event define an event
func (sm *StateMachine) Event(name string) *Event {
	event := &Event{Name: name}
	sm.events[name] = event
	return event
}

// Trigger trigger an event
func (sm *StateMachine) Trigger(name string, value ...interface{}) error {
	stateWas := sm.curState.Name
	if event := sm.events[name]; event != nil {
		var matchedTransitions []*EventTransition
		for _, transition := range event.transitions {
			var validFrom = len(transition.froms) == 0
			if len(transition.froms) > 0 {
				for _, from := range transition.froms {
					if from == stateWas {
						validFrom = true
					}
				}
			}

			if validFrom {
				matchedTransitions = append(matchedTransitions, transition)
			}
		}

		if len(matchedTransitions) == 1 {
			transition := matchedTransitions[0]

			// State: exit
			if state, ok := sm.states[stateWas]; ok {
				if err := state.onExit(); err != nil {
					return err
				}
			}

			if state, ok := sm.states[transition.to]; ok {
				if err := state.onEnter(value...); err != nil {
					return err
				}
			}
			return nil
		}
	}
	return fmt.Errorf("failed to perform event %s from state %s", name, stateWas)
}

// Event contains Event information, including transition hooks
type Event struct {
	Name        string
	transitions []*EventTransition
}

// To define EventTransition of go to a state
func (event *Event) To(name string) *EventTransition {
	transition := &EventTransition{to: name}
	event.transitions = append(event.transitions, transition)
	return transition
}

// EventTransition hold event's to/froms states, also including befores, afters hooks
type EventTransition struct {
	to    string
	froms []string
}

// From used to define from states
func (transition *EventTransition) From(states ...string) *EventTransition {
	transition.froms = states
	return transition
}

type stateFunc func(...interface{}) error
type State struct {
	onEnter stateFunc
	onExit  func() error
	Name    string
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

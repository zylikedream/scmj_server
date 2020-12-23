package tablestate

import (
	"fmt"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/game/table/tableplayer"

	"github.com/pkg/errors"
)

type ITableForState interface {
	UpdateTurnSeat()
	GetTurnPlayer() *tableplayer.TablePlayer
	DrawCard() error
	AfterDiscard() error
	GetPlayers() []*tableplayer.TablePlayer
	GetTurnSeat() int
	SetNextSeat(seat int)
	GetNextTurnPlayer() *tableplayer.TablePlayer
}

const (
	TABLE_STATE_WIN     = "state_win"
	TABLE_STATE_KONG    = "state_kong"
	TABLE_STATE_PONG    = "state_gong"
	TABLE_STATE_DISCARD = "state_discard"
	TABLE_STATE_DRAW    = "state_draw"
)

type StateMachine struct {
	curState IState
	states   map[string]IState
}

func New(table ITableForState) *StateMachine {
	sm := &StateMachine{
		states: map[string]IState{
			TABLE_STATE_DRAW:    NewStateDraw(table),
			TABLE_STATE_DISCARD: NewStateDiscard(table),
		},
	}
	for _, state := range sm.states {
		state.SetStateMachine(sm)
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

	sm.curState = nextState
	if err := nextState.OnEnter(); err != nil {
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

func getStateByOperate(op int) string {
	switch op {
	case tableoperate.OPERATE_WIN:
		return TABLE_STATE_WIN
	case tableoperate.OPERATE_KONG_WIND, tableoperate.OPERATE_KONG_CONCEALED, tableoperate.OPERATE_KONG_EXPOSED:
		return TABLE_STATE_KONG
	case tableoperate.OPERATE_PONG:
		return TABLE_STATE_PONG
	case tableoperate.OPERATE_DISCARD:
		return TABLE_STATE_DISCARD
	default:
		return TABLE_STATE_DRAW
	}
}

type IStateMachine interface {
	GetState(name string) IState
}

type IState interface {
	GetName() string
	OnEnter() error
	OnUpdate() (IState, error)
	OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error
	OnExit() error
	SetStateMachine(sm IStateMachine)
}

type StateBase struct {
	name         string
	stateMachine IStateMachine
}

func (s *StateBase) GetName() string {
	return s.name
}

func (s *StateBase) OnEnter() error {
	return nil
}

func (s *StateBase) OnUpdate() (IState, error) {
	return nil, nil
}

func (s *StateBase) OnExit() error {
	return nil
}

func (s *StateBase) OnPlyOperate(pid uint64, data tableoperate.OperateCommand) error {
	return nil
}

func (s *StateBase) SetStateMachine(sm IStateMachine) {
	s.stateMachine = sm
}

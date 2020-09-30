package tablestate

import (
	"testing"
)

type TestState struct {
	StateBase
}

func (ts *TestState) onEnter(t *testing.T, pid int, card int) error {
	t.Logf("pid=%d, card=%d\n", pid, card)
	return nil
}

func (ts *TestState) OnEnter(f ...interface{}) error {
	return ts.StateBase.OnEnter(f...)
}

func TestStateFunc(t *testing.T) {
	st := &TestState{
		StateBase: StateBase{"test"},
	}
	st.OnEnter(t, 1, 2)
}

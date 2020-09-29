package tablestate

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStateFunc(t *testing.T) {
	st := &state{}
	if err := setEnter(t, st, testFunc); err != nil {
		t.Errorf("set func failed, err=%s", err)
	}
	if err := st.onEnter(1234, 10); err != nil {
		t.Errorf("call on enter failed, err=%s", err)
	}
}

func setEnter(t *testing.T, st *state, f interface{}) error {
	ftype := reflect.TypeOf(f)
	if ftype.Kind() != reflect.Func {
		return fmt.Errorf("need a func")
	}
	st.onEnter = func(args ...interface{}) error {
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
	}
	return nil
}

func testFunc(pid int, card int) error {
	fmt.Printf("pid=%d, card=%d\n", pid, card)
	return nil
}

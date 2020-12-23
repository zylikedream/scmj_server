package util

import (
	"errors"
	"fmt"
	"testing"
)

var errTest1 = errors.New("")
var errTest2 = fmt.Errorf("%w", errTest1)
var errTest3 = fmt.Errorf("%w", errTest2)

func getError() error {
	return errTest3
}

func TestErrorWrap(t *testing.T) {
	err := getError()
	t.Logf("err is errTest1:%t", errors.Is(err, errTest1))
	t.Logf("err is errTest2:%t", errors.Is(err, errTest2))
	t.Logf("err is errTest3:%t", errors.Is(err, errTest3))
}

func TestRemoveElemWithoutOrder(t *testing.T) {
	a := []int{0, 1, 2, 3, 4}
	b := []uint64{0, 1, 2, 3, 4}
	c := []string{"a", "b", "c", "d"}
	item := RemoveElemWithoutOrder(1, &a)
	fmt.Println(a, item)
	RemoveElemWithoutOrder(2, &b)
	fmt.Println(b)
	RemoveElemWithoutOrder(3, &c)
	RemoveElemWithoutOrder(4, &c)
	fmt.Println(c)
}

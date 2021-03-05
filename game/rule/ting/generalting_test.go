package ting

import (
	"testing"
	"zinx-mj/game/rule/win"
)

func Test_generalTing_GetTingCard(t *testing.T) {
	tingRule := NewGeneralRule()
	cards := []int{11, 11, 11, 11, 12, 12, 12, 12, 13, 13, 13, 13, 17}
	tc := tingRule.GetTingCard(cards, win.NewGeneralWin(14))
	t.Logf("tingCard=%#v", tc)
}

package win

import "testing"

func BenchmarkCommonWin_CanWin(b *testing.B) {
	winRule := NewGeneralWin()
	//handCards := []int{6, 7, 9, 9, 12, 12, 13, 14, 15, 15, 17, 26, 27, 28}
	//handCards := []int{ 1, 2, 2, 3, 3, 4, 4, 4, 4, 5, 6, 6}
	//handCards := []int{1, 1, 5, 5, 9, 9, 11, 11, 12, 12, 13, 13, 19, 25}
	handCards := []int{1, 1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 9, 9}
	for i := 0; i < b.N; i++ {
		winRule.CanWin(handCards)
	}
}

func TestCanWin(t *testing.T) {
	winRule := NewGeneralWin()
	var handCards []int
	// 是否可以胡7对
	handCards = []int{1, 1, 5, 5, 9, 9, 11, 11, 12, 12, 13, 13, 19, 19}
	if !winRule.CanWin(handCards) {
		t.Error("7对类型验证失败")
	}
	handCards = []int{1, 1, 5, 5, 9, 9, 11, 11, 12, 12, 13, 13, 19, 25}
	if winRule.CanWin(handCards) {
		t.Error("7对类型验证失败")
	}

	if winRule.CanWin(handCards) {
		t.Error("地龙牌型判断错误")
	}
	handCards = []int{1, 1, 5, 5, 9, 9, 11, 11, 12, 12, 18}
	if winRule.CanWin(handCards) {
		t.Error("地龙牌型判断错误")
	}

	// 是否可以胡单吊
	handCards = []int{1, 1}
	if !winRule.CanWin(handCards) {
		t.Error("单吊牌型验证失败")
	}
	handCards = []int{1, 2}
	if winRule.CanWin(handCards) {
		t.Error("单吊牌型验证失败")
	}

	// 是否可以胡5张手牌
	handCards = []int{1, 1, 1, 2, 3}
	if !winRule.CanWin(handCards) {
		t.Error("5张牌型验证失败")
	}

	// 是否可以胡5张手牌
	handCards = []int{1, 1, 1, 2, 4}
	if winRule.CanWin(handCards) {
		t.Error("5张牌型验证失败")
	}

	// 是否可以胡8张手牌
	handCards = []int{1, 1, 1, 1, 2, 3, 9, 9}
	if !winRule.CanWin(handCards) {
		t.Error("5张牌型验证失败")
	}
	handCards = []int{1, 1, 1, 1, 2, 3, 9, 12}
	if winRule.CanWin(handCards) {
		t.Error("5张牌型验证失败")
	}
	// 是否可以胡11张手牌
	handCards = []int{1, 1, 1, 2, 2, 3, 3, 4, 4, 5, 5}
	if !winRule.CanWin(handCards) {
		t.Error("11张牌型验证失败1")
	}
	// 是否可以胡ABBCCCDDMM
	handCards = []int{1, 2, 2, 3, 3, 3, 4, 4, 5, 7, 7}
	if !winRule.CanWin(handCards) {
		t.Error("11张牌型验证失败1")
	}

	handCards = []int{1, 2, 3, 3, 3, 3, 4, 4}
	if !winRule.CanWin(handCards) {
		t.Error("11张牌型验证失败1")
	}

	handCards = []int{3, 4, 4, 4, 4, 5, 6, 6}
	if !winRule.CanWin(handCards) {
		t.Error("11张牌型验证失败1")
	}

	handCards = []int{1, 2, 2, 3, 3, 4, 4, 4, 4, 5, 6, 6}
	if winRule.CanWin(handCards) {
		t.Error("11张牌型验证失败1")
	}
}

package ting

import (
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"

	"github.com/aceld/zinx/zlog"
)

type generalTing struct {
}

func NewGeneralRule() irule.ITing {
	return &generalTing{}
}

/*
 * Descrp: 得到玩家可以听的牌
 * Create: zhangyi 2020-08-03 22:20:26
 */
func (g *generalTing) GetTingCard(cards []int, winRule irule.IWin) map[int]struct{} {
	tingCards := make(map[int]struct{})
	// 循环将可能听的牌，带入到手牌，再用胡牌算法检测是否可胡
	maybeCards := getMaybeTing(cards)
	defer func() {
		if err := recover(); err != nil {
			zlog.Errorf("ting error: cards:%v maybeCards:%v", cards, maybeCards)
			panic(err)
		}
	}()
	for c := range maybeCards {
		mayWinCards := append([]int{c}, cards...)
		if winRule.CanWin(mayWinCards) {
			tingCards[c] = struct{}{}
		}
	}
	return tingCards
}

func (g *generalTing) CanTing(cards []int, winRule irule.IWin) bool {
	return len(g.GetTingCard(cards, winRule)) > 0
}

/*
 * Descrp: 获取可能的听牌
 * Notice: 一个牌的自身和左右两张牌是可能听的牌
 * Create: zhangyi 2020-08-03 22:08:15
 */
func getMaybeTing(cards []int) map[int]struct{} {
	maybeCards := map[int]struct{}{}
	for _, c := range cards {
		nbor := gamedefine.GetNeighborCards(c)
		for _, n := range nbor {
			maybeCards[n] = struct{}{}
		}
	}
	return maybeCards
}

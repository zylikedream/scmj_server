package ting

import (
	"zinx-mj/game/card"
	"zinx-mj/game/rule/irule"
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
	for _, c := range maybeCards {
		cards = append(cards, c)
		if winRule.CanWin(cards) {
			tingCards[c] = struct{}{}
		}
		cards = cards[:len(cards)-1]
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
func getMaybeTing(cards []int) []int {
	var maybeCards []int
	for _, c := range cards {
		maybeCards = append(maybeCards, card.GetNeighborCards(c)...)
	}
	return maybeCards
}

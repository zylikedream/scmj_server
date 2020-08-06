// 通用发牌规则
package deal

import "zinx-mj/game/rule/irule"

type generalDeal struct {
}

func NewGeneralDeal() irule.IDeal {
	return &generalDeal{}
}

func (g *generalDeal) Deal(cards []int, count int) []int {
	dealCards := cards[:count]
	cards = cards[count:]
	return dealCards
}

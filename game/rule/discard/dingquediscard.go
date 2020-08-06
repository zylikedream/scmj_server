package discard

import (
	"fmt"
	"zinx-mj/game/card"
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/rule/irule"
)

// 定缺出牌规则
type dingqueDisacard struct {
}

func NewDingQueDiscard() irule.IDiscard {
	return &dingqueDisacard{}
}

func (d *dingqueDisacard) Discard(pc *playercard.PlayerCard, crd int, dingqueSuit int) (*playercard.PlayerCard, error) {
	// 如果还有定缺的牌没有打完就必须先打定缺花型的牌
	if card.GetCardSuit(crd) != dingqueSuit && len(pc.GetCardBySuit(dingqueSuit)) > 0 {
		return pc, fmt.Errorf("discard failed, must discard dingque card first")
	}
	if err := pc.Discard(crd); err != nil {
		return pc, err
	}
	return pc, nil
}

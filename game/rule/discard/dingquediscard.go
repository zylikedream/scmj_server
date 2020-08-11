package discard

import (
	"fmt"
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
)

// 定缺出牌规则
type dingqueDisacard struct {
}

func NewDingQueDiscard() irule.IDiscard {
	return &dingqueDisacard{}
}

func (d *dingqueDisacard) Discard(pc *playercard.PlayerCard, crd int, dingqueSuit int) error {
	// 如果还有定缺的牌没有打完就必须先打定缺花型的牌
	if gamedefine.GetCardSuit(crd) != dingqueSuit && len(pc.GetCardBySuit(dingqueSuit)) > 0 {
		return fmt.Errorf("discard failed, must discard dingque card first")
	}
	if err := pc.Discard(crd); err != nil {
		return err
	}
	return nil
}

package discard

import (
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"

	"github.com/pkg/errors"
)

// 定缺出牌规则
type dingqueDisacard struct {
}

func NewDingQueDiscard() irule.IDiscard {
	return &dingqueDisacard{}
}

func (d *dingqueDisacard) Discard(pc *handcard.HandCard, crd int) error {
	// 如果还有定缺的牌没有打完就必须先打定缺花型的牌
	dingqueSuit := pc.DingQueSuit
	if gamedefine.GetCardSuit(crd) != dingqueSuit && len(pc.GetCardBySuit(dingqueSuit)) > 0 {
		return errors.Errorf("discard failed, must discard dingque card first")
	}
	return nil
}

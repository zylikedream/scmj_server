package draw

import (
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/rule/irule"
)

// 通用摸牌规则
type generalDraw struct {
}

func NewGeneralDraw() irule.IDraw {
	return generalDraw{}
}

func (g generalDraw) Draw(pc *handcard.HandCard, card int) error {
	if err := pc.Draw(card); err != nil {
		return err
	}
	return nil
}

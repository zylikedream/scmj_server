package draw

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/rule/irule"
)

// 通用摸牌规则
type generalDraw struct {
}

func NewGeneralDraw() irule.IDraw {
	return generalDraw{}
}

func (g generalDraw) Draw(pc *playercard.PlayerCard, card int) (*playercard.PlayerCard, error) {
	if err := pc.Draw(card); err != nil {
		return pc, err
	}
	return pc, nil
}

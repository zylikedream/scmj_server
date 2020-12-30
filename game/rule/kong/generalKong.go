package kong

import (
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/rule/irule"
)

type generalKong struct {
}

func NewGeneralKong() irule.IKong {
	return &generalKong{}
}

func (g *generalKong) Kong(pc *handcard.HandCard, card int) error {
	if err := pc.Kong(card); err != nil {
		return err
	}
	return nil
}

func (g *generalKong) ExposedKong(pc *handcard.HandCard, card int) error {
	if err := pc.ExposedKong(card); err != nil {
		return err
	}
	return nil
}

func (g *generalKong) ConcealedKong(pc *handcard.HandCard, card int) error {
	if err := pc.ConcealedKong(card); err != nil {
		return err
	}
	return nil
}

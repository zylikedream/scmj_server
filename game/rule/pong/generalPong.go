package pong

import (
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/rule/irule"
)

type generalPong struct {
}

func NewGeneralPong() irule.IPong {
	return &generalPong{}
}

func (g *generalPong) Pong(pc *handcard.HandCard, card int) error {
	if err := pc.Pong(card); err != nil {
		return err
	}
	return nil
}

func (g *generalPong) CanPong(pc *handcard.HandCard, card int) bool {
	return pc.GetCardNum(card) >= 2
}

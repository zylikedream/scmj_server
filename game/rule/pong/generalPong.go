package pong

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/rule/irule"
)

type generalPong struct {
}

func NewGeneralPong() irule.IPong {
	return &generalPong{}
}

func (g *generalPong) Pong(pc *playercard.PlayerCard, card int) error {
	if err := pc.Pong(card); err != nil {
		return err
	}
	return nil
}

func (g *generalPong) CanPong(pc *playercard.PlayerCard, card int) bool {
	return pc.GetCardNum(card) >= 2
}

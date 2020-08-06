package pong

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/rule/irule"
)

type generalPong struct {
}

func NewGeneralKong() irule.IPong {
	return &generalPong{}
}

func (g *generalPong) Pong(pc *playercard.PlayerCard, card int) (*playercard.PlayerCard, error) {
	if err := pc.Pong(card); err != nil {
		return pc, err
	}
	return pc, nil
}

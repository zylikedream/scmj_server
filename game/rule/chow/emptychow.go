package chow

import (
	"errors"
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/rule/irule"
)

// 空吃牌 不允许吃牌
type emptyChow struct {
}

func NewEmptyChow() irule.IChow {
	return &emptyChow{}
}

func (e *emptyChow) Chow(cards *playercard.PlayerCard, card int) error {
	return errors.New("chow card not allowed")
}

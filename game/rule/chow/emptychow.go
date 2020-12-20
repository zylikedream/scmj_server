package chow

import (
	"errors"
	handcard "zinx-mj/game/card/handcard"
	"zinx-mj/game/rule/irule"
)

// 空吃牌 不允许吃牌
type emptyChow struct {
}

func NewEmptyChow() irule.IChow {
	return &emptyChow{}
}

func (e *emptyChow) Chow(cards *handcard.HandCard, card int) error {
	return errors.New("chow card not allowed")
}

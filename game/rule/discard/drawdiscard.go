package discard

import (
	"fmt"
	"zinx-mj/game/card/playercard"
	"zinx-mj/game/rule/irule"
)

// 只能出摸到的牌
type drawDiscard struct {
}

func NewDrawDiscard() irule.IDiscard {
	return drawDiscard{}
}

/*
 * Descrp: 玩家必须摸到哪张打哪张
 * Create: zhangyi 2020-07-03 14:44:52
 */
func (d drawDiscard) Discard(pc *playercard.PlayerCard, crd int, dingqueSuit int) error {
	if crd != pc.GetLastDraw() {
		return fmt.Errorf("can't discard card different with draw card, card=%d, draw=%d",
			crd, pc.GetLastDraw())
	}
	if err := pc.Discard(crd); err != nil {
		return err
	}
	return nil
}

package discard

import (
	"fmt"
	handcard "zinx-mj/game/card/handcard"
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
func (d drawDiscard) Discard(pc *handcard.HandCard, crd int, dingqueSuit int) error {
	if crd != pc.GetLastDraw() {
		return fmt.Errorf("can't discard card different with draw card, card=%d, draw=%d",
			crd, pc.GetLastDraw())
	}
	return nil
}

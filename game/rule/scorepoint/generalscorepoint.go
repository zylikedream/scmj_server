package scorepoint

import (
	"zinx-mj/game/rule/cardmodel"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/winmode"
)

type generalScorePoint struct {
}

func NewGeneralScorePoint() *generalScorePoint {
	return &generalScorePoint{}
}

func (g *generalScorePoint) ScorePoint(models []int, winMode int, kongs int, handcard map[int]int) irule.ScorePoint {
	var scores []irule.ScorePoint
	scores = append(scores, g.scoreCardModelPoint(models))
	scores = append(scores, g.scoreKongPoint(kongs))
	scores = append(scores, g.scoreWinModePoint(winMode))
	scores = append(scores, g.scoreHandCardPoint(handcard))
	total := irule.ScorePoint{}
	for _, s := range scores {
		total.Point += s.Point
		total.Base += s.Base
	}
	return total
}

// 数番
func (g *generalScorePoint) scoreCardModelPoint(models []int) irule.ScorePoint {
	var points int
	for _, model := range models {
		var modelPoint int
		switch model {
		case cardmodel.CARD_MODEL_BAMBOO_SECOND:
			modelPoint = 1
		case cardmodel.CARD_MODEL_BIG_PAIR:
			modelPoint = 2
		case cardmodel.CARD_MODEL_SINGLE_SUIT:
			modelPoint = 2
		case cardmodel.CARD_MODEL_SEVEN_PAIR:
			modelPoint = 2
		case cardmodel.CARD_MODEL_BIG_SEVEN_PAIR: // 龙七对其实还是一种暗七对只有两番，杠另外算番，所以会有三番
			modelPoint = 2
		default:
			modelPoint = 0
		}
		points += modelPoint
	}
	return irule.ScorePoint{
		Point: points,
	}
}

func (g *generalScorePoint) scoreWinModePoint(winMode int) irule.ScorePoint {
	// 这儿自摸当做加底
	score := irule.ScorePoint{}
	switch winMode {
	case winmode.WIN_DRAW_MODE_GOD, winmode.WIN_DISCARD_MODE_DEVIL:
		score.Point += 4 // todo从规则取最大番数
	case winmode.WIN_DRAW_MODE_KONG:
		score.Point += 1
	case winmode.WIN_DISCARD_MODE_KONG, winmode.WIN_DISCARD_MODE_RUB_KONG:
		score.Point += 1
	case winmode.WIN_DISCARD_MODE_PLAIN: // 平胡没有番
	}
	if winmode.IsDrawWin(winMode) {
		score.Base += 1
	}
	return score
}

func (g *generalScorePoint) scoreKongPoint(kongs int) irule.ScorePoint {
	return irule.ScorePoint{
		Point: kongs,
	}
}

func (g *generalScorePoint) scoreHandCardPoint(handcard map[int]int) irule.ScorePoint {
	var handKong int
	// 手牌中的杠
	for _, cardNum := range handcard {
		if cardNum == 4 {
			handKong += 1
		}
	}
	return irule.ScorePoint{
		Point: handKong,
		Base:  0,
	}
}

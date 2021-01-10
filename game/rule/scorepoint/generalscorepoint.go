package scorepoint

import (
	"zinx-mj/game/card/handcard"
	"zinx-mj/game/rule/cardmodel"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/winmode"
)

type generalScorePoint struct {
}

func NewGeneralScorePoint() *generalScorePoint {
	return &generalScorePoint{}
}

func (g *generalScorePoint) ScorePoint(models []int, winMode int, kongs []handcard.KongInfo) irule.ScorePoint {
	cardModelScore := g.scoreCardModelPoint(models)
	kongScore := g.scoreKongPoint(kongs)
	ruleScore := g.scoreWinModePoint(winMode)
	return irule.ScorePoint{
		Point: cardModelScore.Point + kongScore.Point + ruleScore.Point,
		Base:  cardModelScore.Base + kongScore.Base + ruleScore.Base,
	}
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
		case cardmodel.CARD_MODEL_BIG_SEVEN_PAIR:
			modelPoint = 3
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
	case winmode.WIN_SELF_DRAW_MODE_GOD, winmode.WIN_DISCARD_MODE_DEVIL:
		score.Point = 4 // todo从规则取最大番数
	case winmode.WIN_SELF_DRAW_MODE_KONG:
		score.Point = 1
	case winmode.WIN_DISCARD_MODE_KONG, winmode.WIN_DISCARD_MODE_RUB_KONG:
		score.Point = 1
	case winmode.WIN_SELF_DRAW_MODE_PLAIN:
		score.Base = 1 // 自摸加底
	case winmode.WIN_DISCARD_MODE_PLAIN:
		score.Point = 0 // 平胡
	}
	return score
}

func (g *generalScorePoint) scoreKongPoint(kongs []handcard.KongInfo) irule.ScorePoint {
	return irule.ScorePoint{
		Point: len(kongs),
	}
}

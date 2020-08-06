package scmjrule

import (
	"zinx-mj/game/rule/board"
	"zinx-mj/game/rule/chow"
	"zinx-mj/game/rule/deal"
	"zinx-mj/game/rule/discard"
	"zinx-mj/game/rule/draw"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/kong"
	"zinx-mj/game/rule/pong"
	"zinx-mj/game/rule/shuffle"
	"zinx-mj/game/rule/ting"
	"zinx-mj/game/rule/win"
)

type ScmjRule struct {
	boardRule   irule.IBoard
	chowRule    irule.IChow
	discardRule irule.IDiscard
	drawRule    irule.IDraw
	kongRule    irule.IKong
	pongRule    irule.IPong
	shuffleRule irule.IShuffle
	tingRule    irule.ITing
	winRule     irule.IWin
	dealRule    irule.IDeal
}

func (s ScmjRule) GetBoardRule() irule.IBoard {
	return s.boardRule
}

func (s ScmjRule) GetChowRule() irule.IChow {
	return s.chowRule
}

func (s ScmjRule) GetDiscardRule() irule.IDiscard {
	return s.discardRule
}

func (s ScmjRule) GetDrawRule() irule.IDraw {
	return s.drawRule
}

func (s ScmjRule) GetKongRule() irule.IKong {
	return s.kongRule
}

func (s ScmjRule) GetPongRule() irule.IPong {
	return s.pongRule
}

func (s ScmjRule) GetShuffleRule() irule.IShuffle {
	return s.shuffleRule
}

func (s ScmjRule) GetTingRule() irule.ITing {
	return s.tingRule
}

func (s ScmjRule) GetWinRule() irule.IWin {
	return s.winRule
}

func (s ScmjRule) GetDealRule() irule.IDeal {
	return s.dealRule
}

func NewScmjRule() irule.IMjRule {
	return ScmjRule{
		boardRule:   board.NewThreeSuitBoard(),
		chowRule:    chow.NewEmptyChow(),
		discardRule: discard.NewDingQueDiscard(),
		drawRule:    draw.NewGeneralDraw(),
		kongRule:    kong.NewGeneralKong(),
		pongRule:    pong.NewGeneralKong(),
		shuffleRule: shuffle.NewRandomShuffle(),
		tingRule:    ting.NewGeneralRule(),
		winRule:     win.NewGeneralWin(),
		dealRule:    deal.NewGeneralDeal(),
	}
}

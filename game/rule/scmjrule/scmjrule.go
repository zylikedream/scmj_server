package scmjrule

import (
	"errors"
	"zinx-mj/game/card"
	"zinx-mj/game/room/iroom"
	"zinx-mj/game/room/player"
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
	"zinx-mj/util"
)

type ScmjRule struct {
	room           iroom.IMjRoom
	curPlayerIndex int

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

func (s *ScmjRule) GetCurPlayer() *player.RoomPlayer {
	return s.room.GetPlayer(s.curPlayerIndex)
}

func (s *ScmjRule) IsPlayerTurn(pid int) bool {
	return s.GetCurPlayer() == s.room.GetPlayer(pid)
}

func (s *ScmjRule) Chow(pid int, c int) error {
	ply := s.room.GetPlayer(pid)
	return s.chowRule.Chow(ply.PlyCard, c)
}

func (s *ScmjRule) Discard(pid int, c int) error {
	ply := s.room.GetPlayer(pid)
	return s.discardRule.Discard(ply.PlyCard, c, card.CARD_SUIT_EMPTY)
}

func (s *ScmjRule) Draw(pid int, c int) error {
	ply := s.room.GetPlayer(pid)
	return s.drawRule.Draw(ply.PlyCard, c)
}

func (s *ScmjRule) Kong(pid int, c int) error {
	ply := s.room.GetPlayer(pid)
	return s.kongRule.Kong(ply.PlyCard, c)
}

func (s *ScmjRule) Pong(pid int, c int) error {
	ply := s.room.GetPlayer(pid)
	if s.IsPlayerTurn(pid) {
		return errors.New("can't pong in other turn")
	}
	return s.pongRule.Pong(ply.PlyCard, c)
}

func (s *ScmjRule) Shuffle() error {
	//return s.shuffleRule.Shuffle()
	return nil
}

func (s *ScmjRule) UpdateTingCard(pid int) error {
	ply := s.room.GetPlayer(pid)
	ply.PlyCard.TingCard = s.tingRule.GetTingCard(util.IntMapToIntSlice(ply.PlyCard.HandCardMap), s.winRule)
	return nil
}

func (s *ScmjRule) Win(pid int, c int) error {
	ply := s.room.GetPlayer(pid)
	if _, ok := ply.PlyCard.TingCard[c]; !ok {
		return errors.New("can't win not in ting list")
	}
	return nil
}

func (s *ScmjRule) Deal() error {
	return nil
}

func NewScmjRule(room iroom.IMjRoom) irule.IMjRule {
	return &ScmjRule{
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
		room:        room,
	}
}

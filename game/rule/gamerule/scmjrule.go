package gamerule

import (
	"errors"
	"github.com/golang/protobuf/proto"
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
	"zinx-mj/game/table/itable"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
)

type ScmjRuleData struct {
	PlayMode      uint32 // 玩法模式
	GameTurn      uint32 // 游戏轮数
	MaxPoints     uint32 // 最大番数
	SelfWinType   uint32 // 自摸类型（加底或者加番）
	ExposeWinType uint32 // 杠上炮类型（算点炮还是自摸）
	HszSwitch     uint32 // 是否换三张
	JdSwitch      uint32 // 是否算将对
	MqzzSwitch    uint32 // 是否算门清中张
	TdhSwitch     uint32 // 是否算天地胡
}

func (s *ScmjRuleData) PackToPBMsg() proto.Message {
	return &protocol.ScmjRule{
		PlayMode:      s.PlayMode,
		PlayTurn:      s.GameTurn,
		MaxPoint:      s.MaxPoints,
		SelfWinType:   s.SelfWinType,
		ExposeWinType: s.ExposeWinType,
		HszSwitch:     s.HszSwitch,
		JdSwitch:      s.JdSwitch,
		MqzzSwitch:    s.MqzzSwitch,
		TdhSwitch:     s.TdhSwitch,
	}
}

func (s *ScmjRuleData) UnpackFromPBMsg(msg proto.Message) error {
	rule, ok := msg.(*protocol.ScmjRule)
	if !ok {
		return errors.New("wrong message type")
	}
	s.GameTurn = rule.GetPlayTurn()
	s.MaxPoints = rule.GetMaxPoint()
	s.SelfWinType = rule.GetSelfWinType()
	s.ExposeWinType = rule.GetExposeWinType()
	s.HszSwitch = rule.GetHszSwitch()
	s.JdSwitch = rule.GetJdSwitch()
	s.MqzzSwitch = rule.GetMqzzSwitch()
	s.TdhSwitch = rule.GetTdhSwitch()
	s.PlayMode = rule.GetPlayMode()
	return nil
}

type ScmjRule struct {
	data           *ScmjRuleData
	table          itable.ITable
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

func (s *ScmjRule) GetCurPlayer() *tableplayer.TablePlayer {
	return s.table.GetPlayerByIndex(s.curPlayerIndex)
}

func (s *ScmjRule) IsPlayerRound(pid player.PID) bool {
	return s.GetCurPlayer() == s.table.GetPlayerByPid(pid)
}

func (s *ScmjRule) GetRuleData() irule.IRuleData {
	return s.data
}

func (s *ScmjRule) Operate(op irule.IOperate) error {
	return nil
}

func NewScmjRule(ruleData *ScmjRuleData, table itable.ITable) *ScmjRule {
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
		table:       table,
		data:        ruleData,
	}
}

package sccardtable

import (
	"errors"
	"zinx-mj/network/protocol"

	"google.golang.org/protobuf/proto"
)

type ScTableData struct {
	PlayMode      uint32 // 玩法模式
	GameTurn      uint32 // 游戏轮数
	MaxPlayer     int    // 最大游戏人数
	MaxPoints     uint32 // 最大番数
	SelfWinType   uint32 // 自摸类型（加底或者加番）
	ExposeWinType uint32 // 杠上炮类型（算点炮还是自摸）
	HszSwitch     uint32 // 是否换三张
	JdSwitch      uint32 // 是否算将对
	MqzzSwitch    uint32 // 是否算门清中张
	TdhSwitch     uint32 // 是否算天地胡
}

func (s *ScTableData) PackToPBMsg() proto.Message {
	return &protocol.ScmjData{
		PlayMode:      s.PlayMode,
		PlayTurn:      s.GameTurn,
		MaxPoint:      s.MaxPoints,
		SelfWinType:   s.SelfWinType,
		ExposeWinType: s.ExposeWinType,
		HszSwitch:     s.HszSwitch,
		JdSwitch:      s.JdSwitch,
		MqzzSwitch:    s.MqzzSwitch,
		TdhSwitch:     s.TdhSwitch,
		MaxPlayer:     uint32(s.MaxPlayer),
	}
}

func (s *ScTableData) UnpackFromPBMsg(msg proto.Message) error {
	rule, ok := msg.(*protocol.ScmjData)
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
	s.MaxPlayer = int(rule.GetMaxPlayer())
	return nil
}

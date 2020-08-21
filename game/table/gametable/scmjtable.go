package gametable

import (
	"github.com/golang/protobuf/proto"
	"time"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
)

type GameTable struct {
	Id       uint32                     // 桌子ID
	players  []*tableplayer.TablePlayer // 房间的玩家
	startTm  int64
	gameRule irule.IRule
}

func (g *GameTable) Operate(op irule.IOperate) error {
	return g.gameRule.Operate(op)
}

func (g *GameTable) GetRule() irule.IRule {
	return g.gameRule
}

func (g *GameTable) SetRule(rule irule.IRule) {
	g.gameRule = rule
}

func (g *GameTable) Start() error {
	return nil
}

func (g *GameTable) GetPlayerByPid(pid player.PID) *tableplayer.TablePlayer {
	for _, ply := range g.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (g *GameTable) GetPlayerByIndex(index int) *tableplayer.TablePlayer {
	if index < 0 || index >= len(g.players) {
		return nil
	}
	return g.players[index]
}

func (g *GameTable) Join(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error) {
	ply := tableplayer.NewTablePlayer(plyData)
	ply.AddIdentity(identity)
	g.players = append(g.players, ply)
	// todo 广播通知
	return ply, nil
}

func (g *GameTable) Quit(pid player.PID) error {
	panic("implement me")
}

func (g *GameTable) GetStartTime() int64 {
	return g.startTm
}

func (g *GameTable) GetTableNumber() uint32 {
	return g.Id
}

func (g *GameTable) PackToPBMsg() proto.Message {
	reply := &protocol.ScScmjTableInfo{}
	reply.TableId = g.Id
	reply.StartTime = g.GetStartTime()
	reply.Rule = g.gameRule.GetRuleData().PackToPBMsg().(*protocol.ScmjRule)
	for _, ply := range g.players {
		plyData := &protocol.TablePlayerData{
			Pid:         ply.Pid,
			Photo:       0,
			Name:        "",
			Identity:    ply.Identity,
			OnlineState: uint32(ply.OnlineState),
		}
		reply.Players = append(reply.Players, plyData)
	}
	return reply
}

func NewTable(master *tableplayer.TablePlayerData) (*GameTable, error) {
	t := &GameTable{
		startTm: time.Now().Unix(),
	}
	_, err := t.Join(master, gamedefine.TABLE_IDENTIY_MASTER|gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		return t, err
	}
	return t, nil
}

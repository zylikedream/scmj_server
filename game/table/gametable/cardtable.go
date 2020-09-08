package gametable

import (
	"time"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/util"

	"github.com/aceld/zinx/zlog"

	"github.com/golang/protobuf/proto"
)

type CardTable struct {
	id       uint32                     // 桌子ID
	players  []*tableplayer.TablePlayer // 房间的玩家
	startTm  int64
	gameRule irule.IRule
}

func (g *CardTable) GetID() uint32 {
	return g.id
}

func (g *CardTable) Operate(operate irule.IOperate) error {
	panic("implement me")
}

func (g *CardTable) GetRule() irule.IRule {
	return g.gameRule
}

func (g *CardTable) SetRule(rule irule.IRule) {
	g.gameRule = rule
}

func (g *CardTable) Start() error {
	return nil
}

func (g *CardTable) GetPlayer(pid player.PID) *tableplayer.TablePlayer {
	for _, ply := range g.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (g *CardTable) Join(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error) {
	ply := tableplayer.NewTablePlayer(plyData)
	ply.AddIdentity(identity)
	g.players = append(g.players, ply)

	// 通知其它玩家该玩家加入了房间 // 其实应该延迟一帧发送，需要等待其他协议
	msg := &protocol.ScJoinTable{}
	msg.Player = g.PackPlayerData(ply)
	msg.SeatIndex = int32(len(g.players)) - 1
	g.BroadCast(protocol.PROTOID_SC_JOIN_TABLE, msg)

	// 人满了就开游戏
	if g.IsFull() {
		g.Start()
	}

	return ply, nil
}

func (g *CardTable) IsFull() bool {
	return len(g.players) >= g.gameRule.GetMaxPlayer()
}

func (g *CardTable) Quit(pid player.PID) error {
	panic("implement me")
}

func (g *CardTable) GetStartTime() int64 {
	return g.startTm
}

func (g *CardTable) GetTableNumber() uint32 {
	return g.id
}

func (g *CardTable) BroadCast(protoID protocol.PROTOID, msg proto.Message) {
	for _, ply := range g.players {
		if err := util.SendMsg(ply.Pid, protoID, msg); err != nil {
			zlog.Errorf("braodcast to player failed, pid=%d, protoID=%d", ply.Pid, protoID)
		}
	}
}

func (g *CardTable) PackPlayerData(ply *tableplayer.TablePlayer) *protocol.TablePlayerData {
	return &protocol.TablePlayerData{
		Pid:         ply.Pid,
		Photo:       0,
		Name:        ply.Name,
		Identity:    ply.Identity,
		OnlineState: uint32(ply.OnlineState),
	}
}

func (g *CardTable) PackToPBMsg() proto.Message {
	reply := &protocol.ScScmjTableInfo{}
	reply.TableId = g.id
	reply.StartTime = g.GetStartTime()
	reply.Rule = g.gameRule.GetRuleData().PackToPBMsg().(*protocol.ScmjRule)
	for _, ply := range g.players {
		reply.Players = append(reply.Players, g.PackPlayerData(ply))
	}
	return reply
}

func NewTable(tableID uint32, master *tableplayer.TablePlayerData) (*CardTable, error) {
	t := &CardTable{
		id:      tableID,
		startTm: time.Now().Unix(),
	}
	_, err := t.Join(master, gamedefine.TABLE_IDENTIY_MASTER|gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		return t, err
	}
	return t, nil
}

package scmjtable

import (
	"github.com/golang/protobuf/proto"
	"time"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/scmjrule"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
)

type ScmjTable struct {
	Id       uint32                     // 桌子ID
	players  []*tableplayer.TablePlayer // 房间的玩家
	startTm  int64
	gameRule irule.IMjRule
}

func (r *ScmjTable) StartGame() error {
	return nil
}

func (r *ScmjTable) Update(tm int64) error {
	return nil
}

func (r *ScmjTable) Kong(pid player.PID, crd int) error {
	return nil
}

func (r *ScmjTable) Pong(pid player.PID, crd int) error {
	return nil
}

func (r *ScmjTable) Chow(pid player.PID, crd int) error {
	return nil
}

func (r *ScmjTable) Win(pid player.PID, crd int) error {
	return nil
}

func (r *ScmjTable) Draw(pid player.PID) error {
	return nil
}

func (r *ScmjTable) Discard(pid player.PID, card int) error {
	return nil
}

func (r *ScmjTable) GetPlayer(pid player.PID) *tableplayer.TablePlayer {
	for _, ply := range r.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (r *ScmjTable) Join(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error) {
	ply := tableplayer.NewTablePlayer(plyData)
	ply.AddIdentity(identity)
	r.players = append(r.players, ply)
	// todo 广播通知
	return ply, nil
}

func (r *ScmjTable) Quit(pid player.PID) error {
	panic("implement me")
}

func (r *ScmjTable) GetStartTime() int64 {
	return r.startTm
}

func (r *ScmjTable) GetMjRule() irule.IMjRule {
	return r.gameRule
}

func (r *ScmjTable) GetID() uint32 {
	return r.Id
}

func (r *ScmjTable) PackToPBMsg() proto.Message {
	reply := &protocol.ScScmjTableInfo{}
	reply.TableId = r.Id
	reply.StartTime = reply.GetStartTime()
	return reply
}

func NewScmjTable(master *tableplayer.TablePlayerData, rule *scmjrule.ScmjRuleData) (*ScmjTable, error) {
	t := &ScmjTable{
		startTm: time.Now().Unix(),
	}
	t.gameRule = scmjrule.NewScmjRule(rule, t)
	_, err := t.Join(master, gamedefine.TABLE_IDENTIY_MASTER|gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		return t, err
	}
	return t, nil
}

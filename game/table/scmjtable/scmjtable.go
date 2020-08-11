package scmjtable

import (
	"time"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/scmjrule"
	"zinx-mj/game/table"
	"zinx-mj/game/table/itable"
	"zinx-mj/game/table/player"
)

type ScmjTable struct {
	players        []*player.TablePlayer // 房间的玩家
	curPlayerIndex int                   // 当前玩家索引
	gameTurn       uint32                // 游戏轮数
	maxPoints      uint32                // 最大番数
	startTm        int64
	gameRule       irule.IMjRule
}

func (r *ScmjTable) StartGame() error {
	return nil
}

func (r *ScmjTable) Update(tm int64) error {
	panic("implement me")
}

func (r *ScmjTable) Kong(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjTable) Pong(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjTable) Chow(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjTable) Win(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjTable) Draw(pid int) (int, error) {
	panic("implement me")
}

func (r *ScmjTable) Discard(pid int, card int) error {
	panic("implement me")
}

func (r *ScmjTable) GetPlayer(pid int) *player.TablePlayer {
	for _, ply := range r.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (r *ScmjTable) Join(plyData *player.TablePlayerData, identity uint32) (*player.TablePlayer, error) {
	ply := player.NewTablePlayer(plyData, identity)
	r.players = append(r.players, ply)
	// todo 广播通知
	return ply, nil
}

func (r *ScmjTable) Quit(pid int) error {
	panic("implement me")
}

func (r *ScmjTable) GetStartTime() int64 {
	return r.startTm
}

func (r *ScmjTable) GetMjRule() irule.IMjRule {
	return r.gameRule
}

func NewScmjTable(master *player.TablePlayerData, playCount uint32, maxPoints uint32) (itable.IMjTable, error) {
	t := &ScmjTable{
		gameTurn:  playCount,
		maxPoints: maxPoints,
		startTm:   time.Now().Unix(),
	}
	t.gameRule = scmjrule.NewScmjRule(t)
	ply, err := t.Join(master, table.TABLE_IDENTIY_MASTER)
	if err != nil {
		return t, err
	}
	ply.AddIdentity(table.TABLE_IDENTIY_PLAYER) // 房主可能也是玩家
	return t, nil
}

package scmjroom

import (
	"time"
	"zinx-mj/game/room"
	"zinx-mj/game/room/iroom"
	"zinx-mj/game/room/player"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/scmjrule"
)

type ScmjRoom struct {
	players        []*player.RoomPlayer // 房间的玩家
	curPlayerIndex int                  // 当前玩家索引
	gameTurn       uint32               // 游戏轮数
	maxPoints      uint32               // 最大番数
	startTm        int64
	gameRule       irule.IMjRule
}

func (r *ScmjRoom) StartGame() error {
	return nil
}

func (r *ScmjRoom) Update(tm int64) error {
	panic("implement me")
}

func (r *ScmjRoom) Kong(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjRoom) Pong(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjRoom) Chow(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjRoom) Win(pid int, crd int) error {
	panic("implement me")
}

func (r *ScmjRoom) Draw(pid int) (int, error) {
	panic("implement me")
}

func (r *ScmjRoom) Discard(pid int, card int) error {
	panic("implement me")
}

func (r *ScmjRoom) GetCurPlayer() *player.RoomPlayer {
	return r.players[r.curPlayerIndex]
}

func (r *ScmjRoom) GetPlayer(pid int) *player.RoomPlayer {
	for _, ply := range r.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (r *ScmjRoom) Join(plyData *player.RoomPlayerData, identity uint32) (*player.RoomPlayer, error) {
	ply := player.NewRoomPlayer(plyData, identity)
	r.players = append(r.players, ply)
	// todo 广播通知
	return ply, nil
}

func (r *ScmjRoom) Quit(pid int) error {
	panic("implement me")
}

func (r *ScmjRoom) GetRoomStartTime() int64 {
	return r.startTm
}

func (r *ScmjRoom) GetMjType() int {
	return room.ROOM_MJ_SCMJ
}

func (r *ScmjRoom) GetMjRule() irule.IMjRule {
	return r.gameRule
}

func NewScmjRoom(master *player.RoomPlayerData, playCount uint32, maxPoints uint32) (iroom.IMjRoom, error) {
	r := &ScmjRoom{
		gameTurn:  playCount,
		maxPoints: maxPoints,
		startTm:   time.Now().Unix(),
		gameRule:  scmjrule.NewScmjRule(),
	}
	ply, err := r.Join(master, room.ROOM_IDENTIY_MASTER)
	if err != nil {
		return r, err
	}
	ply.AddIdentity(room.ROOM_IDENTIY_PLAYER) // 房主可能也是玩家
	return r, nil
}

package iroom

import (
	"zinx-mj/game/room/player"
	"zinx-mj/game/rule/irule"
)

type IMjRoom interface {
	// 定时更新
	Update(tm int64) error
	// 杠
	Kong(pid int, crd int) error
	// 碰
	Pong(pid int, crd int) error
	// 吃
	Chow(pid int, crd int) error
	// 胡
	Win(pid int, crd int) error
	// 摸牌
	Draw(pid int) (int, error)
	// 出牌
	Discard(pid int, card int) error
	// 得到当前玩家
	GetCurPlayer() *player.RoomPlayer
	GetPlayer(pid int) *player.RoomPlayer
	// 加入房间
	Join(plyData *player.RoomPlayerData, identity uint32) (*player.RoomPlayer, error)
	// 退出房间
	Quit(pid int) error
	// 得到房间开始时间
	GetRoomStartTime() int64
	// 得到麻将类型
	GetMjType() int
	// 得到麻将玩法
	GetMjRule() irule.IMjRule
	StartGame() error
}

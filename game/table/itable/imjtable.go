package itable

import (
	"github.com/golang/protobuf/proto"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/player"
)

type IMjTable interface {
	// 定时更新
	Update(tm int64) error
	// 杠
	Kong(pid player.PID, crd int) error
	// 碰
	Pong(pid player.PID, crd int) error
	// 吃
	Chow(pid player.PID, crd int) error
	// 胡
	Win(pid player.PID, crd int) error
	// 摸牌
	Draw(pid player.PID) error
	// 出牌
	Discard(pid player.PID, card int) error
	// 得到当前玩家
	GetPlayer(pid player.PID) *tableplayer.TablePlayer
	// 加入房间
	Join(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error)
	// 退出房间
	Quit(pid player.PID) error
	// 得到房间开始时间
	GetStartTime() int64
	// 得到麻将类型
	GetMjRule() irule.IMjRule
	StartGame() error
	GetID() uint32
	PackToPBMsg() proto.Message
}

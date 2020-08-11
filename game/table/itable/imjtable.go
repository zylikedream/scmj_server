package itable

import (
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/player"
)

type IMjTable interface {
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
	GetPlayer(pid int) *player.TablePlayer
	// 加入房间
	Join(plyData *player.TablePlayerData, identity uint32) (*player.TablePlayer, error)
	// 退出房间
	Quit(pid int) error
	// 得到房间开始时间
	GetStartTime() int64
	// 得到麻将类型
	GetMjRule() irule.IMjRule
	StartGame() error
}

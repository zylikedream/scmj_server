package itable

import (
	"github.com/golang/protobuf/proto"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/player"
)

type ITable interface {
	GetPlayer(pid player.PID) *tableplayer.TablePlayer
	// 加入间坐姿
	Join(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error)
	// 退出桌子
	Quit(pid player.PID) error
	// 得到桌子开始时间
	GetStartTime() int64

	GetRule() irule.IRule
	SetRule(rule irule.IRule)

	Start() error
	// 桌子编号
	GetTableNumber() uint32

	PackToPBMsg() proto.Message

	Operate(operate_data irule.IOperate) error
}

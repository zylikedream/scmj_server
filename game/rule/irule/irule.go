package irule

import (
	"github.com/golang/protobuf/proto"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/player"
)

type IRule interface {
	GetCurPlayer() *tableplayer.TablePlayer
	IsPlayerRound(pid player.PID) bool
	GetRuleData() IRuleData
	Operate(op IOperate) error
}

type IRuleData interface {
	PackToPBMsg() proto.Message
	UnpackFromPBMsg(message proto.Message) error
}

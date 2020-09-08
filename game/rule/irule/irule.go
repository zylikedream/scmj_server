package irule

import (
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/player"

	"github.com/golang/protobuf/proto"
)

type IRule interface {
	GetCurPlayer() *tableplayer.TablePlayer
	IsPlayerRound(pid player.PID) bool
	GetRuleData() IRuleData
	GetMaxPlayer() int
}

type IRuleData interface {
	PackToPBMsg() proto.Message
	UnpackFromPBMsg(message proto.Message) error
}

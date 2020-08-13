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

	ChangeRound() error
	LockRound(key int) error
	UnLockRound(key int) error
	IsRoundLocked() bool
}

type IRuleData interface {
	PackToPBMsg() proto.Message
	UnpackFromPBMsg(message proto.Message) error
}

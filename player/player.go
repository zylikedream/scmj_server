package player

import (
	"github.com/aceld/zinx/ziface"
	"zinx-mj/database/define"
)

type PID = uint64

type Player struct {
	Pid           PID
	Account       string // 账号
	Name          string // 名字
	CreateTime    int64  // 创建日期
	LastLoginTime int64  // 上次登录时间
	RoomCard      int64  // 房卡数量
	Sex           uint8  // 性别

	Conn ziface.IConnection `bson:"-"`
}

func New() *Player {
	return &Player{
		Sex: define.SEX_MAIL,
	}
}

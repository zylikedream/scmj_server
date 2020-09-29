package player

import (
	"zinx-mj/database/define"

	"github.com/aceld/zinx/ziface"
)

// 玩家PID类型
type PID = uint64

type Player struct {
	Pid           PID    `bson:"pid"`
	Account       string `bson:"account"`         // 账号
	Name          string `bson:"name"`            // 名字
	CreateTime    int64  `bson:"create_time"`     // 创建日期
	LastLoginTime int64  `bson:"last_login_time"` // 上次登录时间
	RoomCard      int64  `bson:"romm_card"`       // 房卡数量
	TableID       uint32 `bson:"table_id"`        // 桌子id
	Sex           uint8  `bson:"sex"`             // 性别

	Conn ziface.IConnection `bson:"-"`
}

func New() *Player {
	return &Player{
		Sex: define.SEX_MAIL,
	}
}

func (p *Player) SetTableID(tableID uint32) {
	p.TableID = tableID
}

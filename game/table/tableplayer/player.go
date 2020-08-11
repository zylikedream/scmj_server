package tableplayer

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/player"
)

// 解耦TablePlayer和Player直接的关系
type TablePlayerData struct {
	Pid player.PID
}

type TablePlayer struct {
	TablePlayerData
	Identity    uint32 // 身份
	OnlineState byte   // 是否在线
	PlyCard     *playercard.PlayerCard
}

func NewTablePlayer(playerData *TablePlayerData) *TablePlayer {
	return &TablePlayer{
		TablePlayerData: *playerData,
	}
}

/*
 * Descrp: 增加玩家身份
 * Create: zhangyi 2020-08-04 01:31:09
 */
func (r *TablePlayer) AddIdentity(identity uint32) uint32 {
	r.Identity |= identity
	return r.Identity
}

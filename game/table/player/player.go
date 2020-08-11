package player

import "zinx-mj/game/card/playercard"

// 解耦RoomPlayer和Player直接的关系
type TablePlayerData struct {
	Pid int
}

type TablePlayer struct {
	TablePlayerData
	Identity uint32 // 身份
	PlyCard  *playercard.PlayerCard
}

func NewTablePlayer(playerData *TablePlayerData, identity uint32) *TablePlayer {
	return &TablePlayer{
		TablePlayerData: *playerData,
		Identity:        identity,
		PlyCard:         playercard.NewPlayerCard(14),
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

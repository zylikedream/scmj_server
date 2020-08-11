package player

import "zinx-mj/game/card/playercard"

// 解耦RoomPlayer和Player直接的关系
type RoomPlayerData struct {
	Pid int
}

type RoomPlayer struct {
	RoomPlayerData
	Identity uint32 // 身份
	PlyCard  *playercard.PlayerCard
}

func NewRoomPlayer(playerData *RoomPlayerData, identity uint32) *RoomPlayer {
	return &RoomPlayer{
		RoomPlayerData: *playerData,
		Identity:       identity,
		PlyCard:        playercard.NewPlayerCard(14),
	}
}

/*
 * Descrp: 增加玩家身份
 * Create: zhangyi 2020-08-04 01:31:09
 */
func (r *RoomPlayer) AddIdentity(identity uint32) uint32 {
	r.Identity |= identity
	return r.Identity
}

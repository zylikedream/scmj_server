package tableplayer

import (
	"zinx-mj/game/card/playercard"
	"zinx-mj/player"
)

// 解耦TablePlayer和Player直接的关系
type TablePlayerData struct {
	Pid  player.PID
	Name string
}

type TablePlayer struct {
	TablePlayerData
	Identity     uint32 // 身份
	OnlineState  byte   // 是否在线
	PlyCard      *playercard.PlayerCard
	validOperate map[int]struct{}
}

func NewTablePlayer(playerData *TablePlayerData) *TablePlayer {
	return &TablePlayer{
		TablePlayerData: *playerData,
		validOperate:    make(map[int]struct{}),
	}
}

/*
 * Descrp: 增加玩家身份
 * Create: zhangyi 2020-08-04 01:31:09
 */
func (t *TablePlayer) AddIdentity(identity uint32) uint32 {
	t.Identity |= identity
	return t.Identity
}

func (t *TablePlayer) AddOperate(op int) {
	t.validOperate[op] = struct{}{}
}

func (t *TablePlayer) IsOperateValid(op int) bool {
	if _, ok := t.validOperate[op]; ok {
		return true
	}
	return false
}

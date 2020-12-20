package tableplayer

import (
	"sort"
	"zinx-mj/game/card/handcard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/player"

	"github.com/pkg/errors"
)

type ITableForPlayer interface {
	GetWinRule() irule.IWin
	GetDiscardRule() irule.IDiscard
}

// 解耦TablePlayer和Player直接的关系
type TablePlayerData struct {
	Pid  player.PID
	Name string
}

type TablePlayer struct {
	TablePlayerData
	Identity     uint32 // 身份
	OnlineState  byte   // 是否在线
	Hcard        *handcard.HandCard
	validOperate map[int]struct{}
	table        ITableForPlayer
}

func NewTablePlayer(playerData *TablePlayerData, table ITableForPlayer) *TablePlayer {
	ply := &TablePlayer{
		TablePlayerData: *playerData,
		validOperate:    make(map[int]struct{}),
		table:           table,
	}
	return ply
}

/*
 * Descrp: 增加玩家身份
 * Create: zhangyi 2020-08-04 01:31:09
 */
func (t *TablePlayer) AddIdentity(identity uint32) uint32 {
	t.Identity |= identity
	return t.Identity
}

func (t *TablePlayer) AddOperate(ops ...int) {
	for _, op := range ops {
		t.validOperate[op] = struct{}{}
	}
}

func (t *TablePlayer) IsOperateValid(op int) bool {
	if _, ok := t.validOperate[op]; ok {
		return true
	}
	return false
}

// 其他人回合的操作
func (t *TablePlayer) GetOperateOnOtherTurn(c int) []int {
	var ops []int
	if t.Hcard.IsTingCard(c) {
		ops = append(ops, tableoperate.OPERATE_WIN)
	}

	if t.Hcard.GetCardNum(c) == 3 {
		ops = append(ops, tableoperate.OPERATE_KONG)
	}
	if t.Hcard.GetCardNum(c) >= 2 {
		ops = append(ops, tableoperate.OPERATE_PONG)
	}
	if len(ops) > 0 {
		ops = append(ops, tableoperate.OPERATE_PASS)
	}
	sort.Ints(ops) // 按照由下级排序
	return ops

}

// 自己回合可以做的操作
func (t *TablePlayer) GetOperateWithSelfTurn() []int {
	var ops []int
	if t.table.GetWinRule().CanWin(t.Hcard.GetCardArray()) {
		ops = append(ops, tableoperate.OPERATE_WIN)
	}
	var CanKang, CanPong bool
	for c, num := range t.Hcard.HandCardMap {
		if num == 4 {
			CanKang = true
		}
		if _, ok := t.Hcard.PongCards[c]; ok {
			CanPong = true
		}
	}
	if CanKang {
		ops = append(ops, tableoperate.OPERATE_KONG)
	}
	if CanPong {
		ops = append(ops, tableoperate.OPERATE_PONG)
	}
	if CanKang || CanPong {
		ops = append(ops, tableoperate.OPERATE_PASS)
	}
	ops = append(ops, tableoperate.OPERATE_DISCARD) // 自己回合可以打牌
	sort.Ints(ops)                                  // 按照优先级排序
	return ops
}

func (t *TablePlayer) DoOperate(opType int, data tableoperate.OperateData) error {
	switch opType {
	case tableoperate.OPERATE_DISCARD:
		return t.discard(opType, data)
	default:
		return errors.Errorf("unsupport operate, op=%d", opType)
	}
}

func (t *TablePlayer) discard(opType int, data tableoperate.OperateData) error {
	err := t.table.GetDiscardRule().Discard(t.Hcard, data.Card, gamedefine.CARD_SUIT_EMPTY)
	if err != nil {
		return err
	}
	if err = t.Hcard.Discard(data.Card); err != nil {
		return err
	}
	return nil
}

func (t *TablePlayer) InitHandCard(cards []int) {
	t.Hcard = handcard.New(cards, len(cards))
}

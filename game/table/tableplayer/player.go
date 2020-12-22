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
	validOperate []int
	operateLog   []tableoperate.OperateCommand // 玩家的操作数据
	table        ITableForPlayer
}

func NewTablePlayer(playerData *TablePlayerData, table ITableForPlayer) *TablePlayer {
	ply := &TablePlayer{
		TablePlayerData: *playerData,
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
	t.validOperate = append(t.validOperate, ops...)
}

func (t *TablePlayer) GetOperates() []int {
	return t.validOperate
}

func (t *TablePlayer) ClearOperates() {
	t.validOperate = t.validOperate[0:0]
}

func (t *TablePlayer) clearOperate(op int) {
	for i, vop := range t.validOperate {
		if vop == op {
			t.validOperate[i] = t.validOperate[len(t.validOperate)-1]
			t.validOperate = t.validOperate[:len(t.validOperate)-1]
			break
		}
	}
}

func (t *TablePlayer) IsOperateValid(op int) bool {
	for _, vop := range t.validOperate {
		if vop == op {
			return true
		}
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
// 自己回合操作没有跳过选项, 必须要做出操作
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
	ops = append(ops, tableoperate.OPERATE_DISCARD) // 自己回合可以打牌
	sort.Ints(ops)                                  // 按照优先级排序
	return ops
}

func (t *TablePlayer) AddOperateLog(cmd tableoperate.OperateCommand) {
	t.operateLog = append(t.operateLog, cmd) // 记录命令
}

func (t *TablePlayer) DoOperate(cmd tableoperate.OperateCommand) error {
	var err error
	switch cmd.OpType {
	case tableoperate.OPERATE_DISCARD:
		err = t.discard(cmd.OpType, cmd.OpData)
	default:
		return errors.Errorf("unsupport operate, op=%d", cmd.OpType)
	}
	t.clearOperate(cmd.OpType) // 防止重复操作
	t.AddOperateLog(cmd)       // 记录命令
	return err
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

func (t *TablePlayer) DrawCard(c int) error {
	if err := t.Hcard.Draw(c); err != nil {
		return err
	}
	t.AddOperate(t.GetOperateWithSelfTurn()...)
	return nil
}

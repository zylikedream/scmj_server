package tableplayer

import (
	"sort"
	"zinx-mj/game/card/handcard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/player"
	"zinx-mj/util"

	"github.com/aceld/zinx/zlog"
	"github.com/pkg/errors"
)

type ITableForPlayer interface {
	GetWinRule() irule.IWin
	GetDiscardRule() irule.IDiscard
	GetTingRule() irule.ITing
	GetTurnPlayer() *TablePlayer
	GetScoreModelRule() irule.IScoreCardModel
	GetScorePointRule() irule.IScorePoint
	GetWinMode(pid uint64) int
}

// 解耦TablePlayer和Player直接的关系
type TablePlayerData struct {
	Pid  player.PID
	Name string
}

type WinInfo struct {
	Pid     uint64
	Card    int
	Models  []int
	WinMode int
	Score   irule.ScorePoint // 赢的分数
}

type TablePlayer struct {
	TablePlayerData
	Identity     uint32     // 身份
	OnlineState  byte       // 是否在线
	Ready        bool       // 是否准备
	Wins         []*WinInfo // 胡的牌
	Hcard        *handcard.HandCard
	validOperate []int
	operateLog   []tableoperate.OperateCommand // 玩家的操作数据
	table        ITableForPlayer
	Points       []int
}

func NewTablePlayer(playerData *TablePlayerData, table ITableForPlayer) *TablePlayer {
	ply := &TablePlayer{
		TablePlayerData: *playerData,
		table:           table,
		operateLog: []tableoperate.OperateCommand{
			{OpType: tableoperate.OPERATE_EMPTY}, // 哨兵命令
		},
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
	zlog.Debugf("ply add ops, pid:%d, addop:%v, all：%v", t.Pid, ops, t.validOperate)
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
			util.RemoveElemWithoutOrder(i, &t.validOperate)
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

// 出牌后的操作
func (t *TablePlayer) GetOperateWithDiscard(c int) []int {
	var ops []int
	if t.Hcard.IsTingCard(c) {
		ops = append(ops, tableoperate.OPERATE_WIN)
	}

	if t.Hcard.GetCardNum(c) == 3 {
		ops = append(ops, tableoperate.OPERATE_KONG_RAIN)
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

// 摸牌可以做的操作
// 自己回合操作没有跳过选项, 必须要做出操作
func (t *TablePlayer) GetOperateWithDraw() []int {
	var ops []int
	if t.table.GetWinRule().CanWin(t.Hcard.GetHandCard()) {
		ops = append(ops, tableoperate.OPERATE_WIN)
	}
	for c, num := range t.Hcard.CardMap {
		if num == 4 { // 暗杠
			ops = append(ops, tableoperate.OPERATE_KONG_CONCEALED)
			break
		} else if _, ok := t.Hcard.PongCards[c]; ok { // 明杠
			ops = append(ops, tableoperate.OPERATE_KONG_EXPOSED)
			break
		}
	}
	if len(ops) > 0 {
		ops = append(ops, tableoperate.OPERATE_PASS)
	}
	ops = append(ops, tableoperate.OPERATE_DISCARD) // 自己回合可以打牌
	sort.Ints(ops)                                  // 按照优先级排序
	return ops
}

// 其他人明杠可以做的操作
func (t *TablePlayer) GetOperateWithConcealedKong(c int) []int {
	var ops []int
	if t.Hcard.IsTingCard(c) {
		ops = append(ops, tableoperate.OPERATE_WIN)
	}
	if len(ops) > 0 {
		ops = append(ops, tableoperate.OPERATE_PASS)
	}
	return ops
}

func (t *TablePlayer) AddOperateLog(cmd tableoperate.OperateCommand) {
	t.operateLog = append(t.operateLog, cmd) // 记录命令
}

func (t *TablePlayer) GetOperateLog() []tableoperate.OperateCommand {
	return t.operateLog
}

func (t *TablePlayer) DoOperate(cmd tableoperate.OperateCommand) error {
	var err error
	switch cmd.OpType {
	case tableoperate.OPERATE_DISCARD:
		err = t.discard(cmd.OpType, cmd.OpData)
	case tableoperate.OPERATE_WIN:
		err = t.win(cmd.OpType, cmd.OpData)
	case tableoperate.OPERATE_KONG_CONCEALED:
	case tableoperate.OPERATE_KONG_EXPOSED:
	case tableoperate.OPERATE_KONG_RAIN:
	case tableoperate.OPERATE_PONG:
	case tableoperate.OPERATE_PASS:
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
	// 更新玩家听牌
	defer func() {
		if err := recover(); err != nil {
			zlog.Errorf("get ting card error, Hcard:%v", *t.Hcard)
			panic(err)
		}
	}()
	tingRule := t.table.GetTingRule()
	t.Hcard.TingCard = tingRule.GetTingCard(t.Hcard.GetHandCard(), t.table.GetWinRule())

	return nil
}

func (t *TablePlayer) win(opType int, data tableoperate.OperateData) error {
	winInfo := &WinInfo{
		Card: data.Card,
	}
	t.Wins = append(t.Wins, winInfo)
	turnPly := t.table.GetTurnPlayer()
	turnPid := turnPly.Pid
	winInfo.Pid = turnPid

	mcards := &irule.CardModel{
		PongCard: t.Hcard.PongCards,
		KongCard: t.Hcard.KongCards,
		WinCard:  data.Card,
		WiRule:   t.table.GetWinRule(),
	}
	handCardCopy := util.CopyIntMap(t.Hcard.CardMap)
	if turnPid != t.Pid { // 别人点炮
		handCardCopy[data.Card] += 1 // 加上点炮的牌
	}
	mcards.HandCard = handCardCopy
	scoreModelRule := t.table.GetScoreModelRule()
	models := scoreModelRule.ScoreCardModel(mcards)
	winMode := t.table.GetWinMode(t.Pid)
	// 计算番数，番数由三部分组成，一是基本牌型（七对，清一色等），二是番数规则（自摸，天胡，杠上花等） , 三是杠的番数
	scorePointRule := t.table.GetScorePointRule()
	score := scorePointRule.ScorePoint(models, winMode, t.Hcard.KongCards) // 牌型番数
	// 规则
	winInfo.Models = models
	winInfo.WinMode = winMode
	winInfo.Score = score

	return nil
}

func (t *TablePlayer) kongConcealed(opType int, data tableoperate.OperateData) error {
	return nil
}

func (t *TablePlayer) IsWinned() bool {
	return len(t.Wins) > 0
}

func (t *TablePlayer) InitHandCard(cards []int, cardMax int) error {
	t.Hcard = handcard.New(cardMax)
	return t.Hcard.SetHandCard(cards)
}

func (t *TablePlayer) DrawCard(c int) error {
	if err := t.Hcard.Draw(c); err != nil {
		return err
	}
	t.AddOperate(t.GetOperateWithDraw()...)
	return nil
}

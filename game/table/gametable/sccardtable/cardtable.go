package sccardtable

import (
	"math/rand"
	"sort"
	"time"
	"zinx-mj/game/card/boardcard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/board"
	"zinx-mj/game/rule/chow"
	"zinx-mj/game/rule/deal"
	"zinx-mj/game/rule/discard"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/kong"
	"zinx-mj/game/rule/pong"
	"zinx-mj/game/rule/shuffle"
	"zinx-mj/game/rule/ting"
	"zinx-mj/game/rule/win"
	"zinx-mj/game/table/tableoperate"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/game/table/tablestate"
	"zinx-mj/mjerror"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/util"

	"github.com/aceld/zinx/zlog"

	"github.com/pkg/errors"

	"google.golang.org/protobuf/proto"
)

type OperateTimer struct {
	T       *time.Timer
	Pid     uint64
	Operate tableoperate.OperateCommand
}

type ScCardTable struct {
	id           uint32                     // 桌子ID
	players      []*tableplayer.TablePlayer // 房间的玩家
	createTime   time.Time
	startTime    time.Time
	games        int // 游戏局数
	turnSeat     int // 当前回合的玩家
	nextTurnSeat int // 当前操作的玩家

	data         *ScTableData
	events       *ScTableEvent
	board        *boardcard.BoardCard
	stateMachine *tablestate.StateMachine
	opTimers     []*OperateTimer

	boardRule   irule.IBoard
	chowRule    irule.IChow
	discardRule irule.IDiscard
	kongRule    irule.IKong
	pongRule    irule.IPong
	shuffleRule irule.IShuffle
	tingRule    irule.ITing
	winRule     irule.IWin
	dealRule    irule.IDeal
}

func NewTable(tableID uint32, master *tableplayer.TablePlayerData, data *ScTableData) (*ScCardTable, error) {
	t := &ScCardTable{
		id:         tableID,
		createTime: time.Now(),
		data:       data,
	}
	t.boardRule = board.NewThreeSuitBoard() // 三坊
	t.chowRule = chow.NewEmptyChow()
	t.discardRule = discard.NewDingQueDiscard()
	t.kongRule = kong.NewGeneralKong()
	t.pongRule = pong.NewGeneralPong()
	t.shuffleRule = shuffle.NewRandomShuffle()
	t.tingRule = ting.NewGeneralRule()
	t.winRule = win.NewGeneralWin()
	t.dealRule = deal.NewGeneralDeal()

	t.initEvent()
	t.initStateMachine()
	_, err := t.PlayerJoin(master, gamedefine.TABLE_IDENTIY_MASTER|gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (s *ScCardTable) initStateMachine() {
	s.stateMachine = tablestate.New(s)
}

func (s *ScCardTable) initEvent() {
	s.events = NewScTableEvent()
	_ = s.events.Register(EVENT_JOIN, s.onJoinEvent)
	_ = s.events.Register(EVENT_GAME_START, s.onGameStart)
}

func (s *ScCardTable) GetID() uint32 {
	return s.id
}

func (s *ScCardTable) GetPlayerSeat(pid uint64) int {
	for i := range s.players {
		ply := s.players[i]
		if ply.Pid == pid {
			return i
		}
	}
	return len(s.players)
}

func (s *ScCardTable) onGameStart() error {
	// 初始化玩家手牌
	s.games++ // 增加游戏局数

	// 初始化庄家
	s.nextTurnSeat = rand.Intn(len(s.players)) // 随机庄家
	s.UpdateTurnSeat()

	// 稍后广播玩家手牌
	msg := &protocol.ScGameTurnStart{}
	// todo 抽象筛子点数rule
	msg.DiePoint = rand.Int31n(6) + 1

	// 设置初始化状态
	if err := s.stateMachine.SetInitState(tablestate.TABLE_STATE_INIT); err != nil {
		return err
	}

	s.startTime = time.Now()

	s.initializeHandCard()
	// 广播游戏开始消息
	if err := s.broadCast(protocol.PROTOID_SC_GAME_TURN_START, msg); err != nil {
		zlog.Errorf("broadCast game start failed, err=%s", err)
		return err
	}

	// 广播玩家手牌
	err := s.broadCastRaw(protocol.PROTOID_SC_CARD_INFO, func(pid uint64, seat int) proto.Message {
		return s.PackCardInfoForPlayer(pid)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ScCardTable) GetPlayerByPid(pid player.PID) *tableplayer.TablePlayer {
	for _, ply := range s.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (s *ScCardTable) GetPlayerBySeat(seat int) *tableplayer.TablePlayer {
	if seat >= len(s.players) || seat < 0 {
		return nil
	}
	return s.players[seat]
}

func (s *ScCardTable) PlayerJoin(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error) {
	ply := tableplayer.NewTablePlayer(plyData, s)
	ply.AddIdentity(identity)
	s.players = append(s.players, ply)

	s.events.Add(EVENT_JOIN, ply.Pid) // 延迟触发，等待玩家加入房间后再通知

	tableData := s.PackToPBMsg()
	if err := util.SendMsg(plyData.Pid, protocol.PROTOID_SC_TABLE_INFO, tableData); err != nil {
		return nil, err
	}
	return ply, nil
}

func (s *ScCardTable) IsFull() bool {
	return len(s.players) >= s.data.MaxPlayer
}

func (s *ScCardTable) Quit(pid player.PID) error {
	panic("implement me")
}

func (s *ScCardTable) IsGameStart() bool {
	return !s.startTime.IsZero()
}

func (s *ScCardTable) GetCreateTime() time.Time {
	return s.createTime
}

func (s *ScCardTable) GetTableNumber() uint32 {
	return s.id
}

// 广播同样的消息给所有玩家
func (s *ScCardTable) broadCast(protoID protocol.PROTOID, msg proto.Message) error {
	return s.broadCastRaw(protoID, func(pid uint64, seat int) proto.Message {
		return msg
	})
}

func (s *ScCardTable) broadCastRaw(protoID protocol.PROTOID,
	msgGenFunc func(pid uint64, seat int) proto.Message) error {
	for i, ply := range s.players {
		msg := msgGenFunc(ply.Pid, i)
		if msg == nil { // 不发送给玩家
			continue
		}
		if err := util.SendMsg(ply.Pid, protoID, msg); err != nil {
			return errors.Errorf("braodcast to player failed, pid=%d, protoID=%d", ply.Pid, protoID)
		}
	}
	return nil
}

func (s *ScCardTable) PackPlayerData(ply *tableplayer.TablePlayer) *protocol.TablePlayerData {
	return &protocol.TablePlayerData{
		Pid:         ply.Pid,
		Photo:       0,
		Name:        ply.Name,
		Identity:    ply.Identity,
		OnlineState: uint32(ply.OnlineState),
	}
}

func (s *ScCardTable) PackToPBMsg() *protocol.ScScmjTableInfo {
	reply := &protocol.ScScmjTableInfo{}
	reply.TableId = s.id
	reply.StartTime = s.startTime.Unix()
	reply.Data = s.data.PackToPBMsg().(*protocol.ScmjData)
	for _, ply := range s.players {
		reply.Players = append(reply.Players, s.PackPlayerData(ply))
	}
	return reply
}

func (s *ScCardTable) onJoinEvent(pid player.PID) error {
	// 通知其它玩家该玩家加入了房间
	msg := &protocol.ScJoinTable{}
	ply := s.GetPlayerByPid(pid)
	if ply == nil {
		return errors.WithStack(mjerror.ErrPlyNotFound)
	}
	msg.Player = s.PackPlayerData(ply)
	msg.SeatIndex = int32(len(s.players)) - 1
	if err := s.broadCast(protocol.PROTOID_SC_JOIN_TABLE, msg); err != nil {
		return err
	}

	return nil
}

const (
	HAND_CARD_NUM     = 13
	MAX_HAND_CARD_NUM = 14
)

func (s *ScCardTable) initializeHandCard() {
	s.board = s.boardRule.NewBoard()
	s.shuffleRule.Shuffle(s.board.Cards)

	for i := 0; i < s.data.MaxPlayer; i++ {
		_ = s.players[i].InitHandCard(s.board.Cards[:HAND_CARD_NUM], MAX_HAND_CARD_NUM)
		s.board.Cards = s.board.Cards[HAND_CARD_NUM:]
	}
}

func (s *ScCardTable) PackCardInfoForPlayer(pid uint64) *protocol.ScCardInfo {
	msg := &protocol.ScCardInfo{}
	msg.TableCard = &protocol.TableCardData{}
	msg.TableCard.Total = int32(len(s.board.CardsTotal))
	msg.TableCard.Left = int32(len(s.board.Cards))

	for _, ply := range s.players {
		var cards []int
		if ply.Pid == pid {
			cards = ply.Hcard.GetHandCard()
		} else {
			cards = ply.Hcard.GetGuardHandCard()
		}
		handCards := &protocol.HandCardData{}
		for _, card := range cards {
			handCards.Card = append(handCards.Card, int32(card))
		}
		msg.HandCard = append(msg.HandCard, handCards)
	}

	return msg
}

func (s *ScCardTable) Update(delta time.Duration) {
	// 更新玩家operate的timer
	s.events.FireAll()
	s.UpdateOperateTimer()

	// 更新状态机
	if err := s.stateMachine.Update(); err != nil {
		zlog.Errorf("update state machine failed, err:%s", err)
	}
	// 处理所有的事件
}

func (s *ScCardTable) DrawCard() error {
	turnPly := s.GetTurnPlayer()
	card, err := s.board.DrawForward()
	if err != nil {
		return err
	}
	zlog.Infof("player draw card, pid:%d, card:%d", turnPly.Pid, card)
	if err = turnPly.DrawCard(card); err != nil {
		return err
	}
	err = s.broadCastRaw(protocol.PROTOID_SC_DRAW_CARD, func(pid uint64, seat int) proto.Message {
		msg := &protocol.ScDrawCard{Pid: pid, Card: -1}
		if turnPly.Pid == pid { // 摸到的牌只会发给本人
			msg.Card = int32(card)
		}
		return msg
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *ScCardTable) NotifyPlyOperate(ply *tableplayer.TablePlayer) error {
	// todo
	return nil
}

func (s *ScCardTable) GetPlayers() []*tableplayer.TablePlayer {
	return s.players
}

func (s *ScCardTable) GetWinRule() irule.IWin {
	return s.winRule
}

func (s *ScCardTable) GetDiscardRule() irule.IDiscard {
	return s.discardRule
}

func (s *ScCardTable) GetTingRule() irule.ITing {
	return s.tingRule
}

func (s *ScCardTable) OnPlyOperate(pid uint64, operate tableoperate.OperateCommand) error {
	zlog.Infof("on ply operate, pid:%d, operate:%v", pid, operate)
	ply := s.GetPlayerByPid(pid)
	if ply == nil {
		return errors.Errorf("not found player, pid=%d", pid)
	}
	if !ply.IsOperateValid(operate.OpType) {
		return errors.Errorf("unalid op for ply op=%d, pid=%d", operate.OpType, pid)
	}

	if err := s.stateMachine.GetCurState().OnPlyOperate(pid, operate); err != nil {
		return err
	}

	if err := ply.DoOperate(operate); err != nil {
		return err
	}

	// 更新下一个回合的玩家
	if operate.OpType != tableoperate.OPERATE_PASS {
		turnSeat := s.GetTurnSeat()
		maxPlayer := s.data.MaxPlayer
		plySeat := s.GetPlayerSeat(ply.Pid)
		// 靠后的玩家才会更新
		if util.SeatRelative(plySeat, turnSeat, maxPlayer) > util.SeatRelative(s.nextTurnSeat, turnSeat, maxPlayer) {
			s.SetNextSeat(plySeat)
		}
	}

	if err := s.AfterPlyOperate(pid, operate); err != nil {
		return err
	}

	pbOperate := &protocol.ScNotifyOperate{
		Pid:    pid,
		OpType: int32(operate.OpType),
		Data: &protocol.OperateData{
			Card: int32(operate.OpData.Card),
		},
	}
	if err := s.broadCast(protocol.PROTOID_SC_NOTIFY_OPERATE, pbOperate); err != nil {
		return err
	}

	cards := ply.Hcard.GetHandCard()
	sort.Ints(cards)
	zlog.Infof("after operate, card:%v, op:%v", cards, operate)

	return nil
}

func (s *ScCardTable) AfterPlyOperate(pid uint64, operate tableoperate.OperateCommand) error {
	var err error
	switch operate.OpType {
	case tableoperate.OPERATE_DISCARD:
		err = s.AfterDiscard(pid, operate.OpData.Card)
	case tableoperate.OPERATE_KONG_CONCEALED: // 玩家可能会抢杠
		err = s.AfterConcealedKong(pid, operate.OpData.Card)
	}
	// 添加定时器
	s.AddOperateTimer(pid, operate)
	return err
}

func (s *ScCardTable) AddOperateTimer(pid uint64, operate tableoperate.OperateCommand) time.Duration {
	for _, t := range s.opTimers {
		if t.Pid == pid {
			zlog.Errorf("add timer repeated, pid:%d, timers:%v", pid, s.opTimers)
			return -1
		}
	}

	timeout := s.GetOperateTimeout(pid, operate)
	if timeout == 0 { // 0 表示永不超时
		// 添加一个占位元素
		s.opTimers = append(s.opTimers, &OperateTimer{Pid: pid})
		return 0
	}
	tr := time.NewTimer(timeout)
	s.opTimers = append(s.opTimers, &OperateTimer{
		T:       tr,
		Pid:     pid,
		Operate: operate,
	})
	return timeout
}

func (s *ScCardTable) CancelOperateTimer(pid uint64) *time.Timer {
	for i, t := range s.opTimers {
		if t.Pid == pid {
			t.T.Stop()
			util.RemoveElemWithoutOrder(i, &s.opTimers) // 删除这个玩家的timer
			return t.T
		}
	}
	return nil
}

func (s *ScCardTable) GetTimeoutOperate(pid uint64, opType int) tableoperate.OperateCommand {
	ply := s.GetPlayerByPid(pid)
	var opCmd tableoperate.OperateCommand
	switch opType {
	case tableoperate.OPERATE_DISCARD:
		opCmd = tableoperate.OperateCommand{
			OpType: tableoperate.OPERATE_DISCARD,
			OpData: tableoperate.OperateData{
				Card: ply.Hcard.GetRecommandCard(),
			},
		}
	default:
		opCmd = tableoperate.OperateCommand{
			OpType: tableoperate.OPERATE_PASS,
		}
	}
	return opCmd
}

func (s *ScCardTable) UpdateOperateTimer() {
	for _, t := range s.opTimers {
		if t.T == nil { // 不超时的玩家
			continue
		}
		if t.T.Stop() { // 定时器已经被关闭了
			continue
		}
		select {
		case <-t.T.C: // 玩家操作超时
			// 强制替玩家操作
			if err := s.OnPlyOperate(t.Pid, s.GetTimeoutOperate(t.Pid, t.Operate.OpType)); err != nil {
				zlog.Errorf("ply operate failed, pid:%d, err:%s", t.Pid, err)
			}
		default:
		}
	}
}

func (s *ScCardTable) GetOperateTimeout(pid uint64, operate tableoperate.OperateCommand) time.Duration {
	return 10 * time.Second
}

func (s *ScCardTable) GetTurnSeat() int {
	return s.turnSeat
}

func (s *ScCardTable) GetTurnPlayer() *tableplayer.TablePlayer {
	return s.GetPlayerBySeat(s.turnSeat)
}

func (s *ScCardTable) GetNextTurnPlayer() *tableplayer.TablePlayer {
	return s.GetPlayerBySeat(s.nextTurnSeat)
}

// 更新的玩家回合时需要清掉玩家的所有操作
func (s *ScCardTable) UpdateTurnSeat() {
	s.turnSeat = s.nextTurnSeat
	for _, ply := range s.players {
		ply.ClearOperates()
	}
}

func (s *ScCardTable) SetNextSeat(seat int) {
	if seat >= s.data.MaxPlayer {
		seat = 0
	}
	s.nextTurnSeat = seat
}

func (s *ScCardTable) AfterDiscard(pid uint64, c int) error {
	for i := range s.players {
		ply := s.players[i]
		if ply.Pid == pid {
			continue
		}
		ops := ply.GetOperateWithDiscard(c)
		ply.AddOperate(ops...)
	}
	return nil
}

func (s *ScCardTable) AfterConcealedKong(pid uint64, c int) error {
	for i := range s.players {
		ply := s.players[i]
		if ply.Pid == pid {
			continue
		}
		ops := ply.GetOperateWithConcealedKong(c)
		ply.AddOperate(ops...)
	}
	return nil
}

func (s *ScCardTable) SetReady(pid uint64, ready bool) {
	if s.IsGameStart() {
		zlog.Warnf("can't set ready when game start")
		return
	}
	ply := s.GetPlayerByPid(pid)
	if ready == ply.Ready {
		return
	}
	ply.Ready = ready
	readData := &protocol.ScPlayerReady{
		Pid:   pid,
		Ready: ready,
	}
	_ = s.broadCast(protocol.PROTOID_SC_PLAYER_READY, readData)

	if !ready {
		return
	}
	// 检查游戏是否可以开始
	if !s.IsFull() {
		return
	}

	for _, ply := range s.players {
		if !ply.Ready {
			return
		}
	}
	// 触发游戏开始
	s.events.Add(EVENT_GAME_START)
}

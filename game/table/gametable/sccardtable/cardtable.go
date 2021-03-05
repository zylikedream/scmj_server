package sccardtable

import (
	"math/rand"
	"time"
	"zinx-mj/game/card/boardcard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/board"
	"zinx-mj/game/rule/chow"
	"zinx-mj/game/rule/deal"
	"zinx-mj/game/rule/discard"
	"zinx-mj/game/rule/gamemode"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/rule/pong"
	"zinx-mj/game/rule/scorecardmodel"
	"zinx-mj/game/rule/scorepoint"
	"zinx-mj/game/rule/shuffle"
	"zinx-mj/game/rule/ting"
	"zinx-mj/game/rule/win"
	"zinx-mj/game/rule/winmode"
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

type ScCardTable struct {
	id            uint32                     // 桌子ID
	players       []*tableplayer.TablePlayer // 房间的玩家
	createTime    time.Time
	gameStartTime time.Time
	gameEndTime   time.Time
	tableEndTime  time.Time
	games         int    // 游戏局数
	turnSeat      int    // 当前回合的玩家
	nextTurnSeat  int    // 当前操作的玩家
	dealer        uint64 // 庄家

	data         *ScTableData
	events       *ScTableEvent
	board        *boardcard.BoardCard
	stateMachine *tablestate.StateMachine
	discards     []int // 桌子上打出的牌
	secondTicker *time.Ticker

	boardRule      irule.IBoard
	chowRule       irule.IChow
	discardRule    irule.IDiscard
	pongRule       irule.IPong
	shuffleRule    irule.IShuffle
	tingRule       irule.ITing
	dealRule       irule.IDeal
	scoreCardModel irule.IScoreCardModel
	scorePointRule irule.IScorePoint
	winModeModel   irule.IWinModeModel
	gameModeRule   irule.IGameMode
}

func NewTable(tableID uint32, master *tableplayer.TablePlayerData, data *ScTableData) (*ScCardTable, error) {
	t := &ScCardTable{
		id:         tableID,
		createTime: time.Now(),
		data:       data,
	}
	t.boardRule = board.NewSuitBoard(gamedefine.CARD_SUIT_BAMBOO, gamedefine.CARD_SUIT_DOT, gamedefine.CARD_SUIT_CHARACTER) // 三坊
	// t.boardRule = board.NewSuitBoard(gamedefine.CARD_SUIT_BAMBOO, gamedefine.CARD_SUIT_DOT) // 两坊
	// t.boardRule = board.NewSuitBoard(gamedefine.CARD_SUIT_BAMBOO) // 一坊
	t.chowRule = chow.NewEmptyChow()
	t.discardRule = discard.NewDingQueDiscard()
	t.pongRule = pong.NewGeneralPong()
	t.shuffleRule = shuffle.NewRandomShuffle()
	// t.shuffleRule = shuffle.NewSortShuffle()
	t.tingRule = ting.NewGeneralRule()
	t.dealRule = deal.NewGeneralDeal()
	t.scoreCardModel = scorecardmodel.NewGeneralScoreCardModel()
	t.scorePointRule = scorepoint.NewGeneralScorePoint()
	t.winModeModel = winmode.NewWinModeModel()
	t.gameModeRule = gamemode.NewXzddGameMode()

	t.secondTicker = time.NewTicker(time.Second)
	t.initEvent()
	t.initStateMachine()
	_, err := t.PlayerJoin(master, gamedefine.TABLE_IDENTIY_OWNER|gamedefine.TABLE_IDENTIY_PLAYER)
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

func (s *ScCardTable) initMaster() {
	s.nextTurnSeat = rand.Intn(len(s.players)) // 随机庄家
	master := s.GetPlayerBySeat(s.nextTurnSeat)
	master.AddIdentity(gamedefine.TABLE_IDENTIY_DEALER)
	s.dealer = master.Pid
}

func (s *ScCardTable) onGameStart() error {
	if s.IsEnd() {
		zlog.Errorf("game start failed, table is end")
		return nil
	}
	// 初始化玩家手牌
	s.games++ // 增加游戏局数

	// 初始化庄家
	s.initMaster()
	s.UpdateTurnSeat()

	// 设置初始化状态
	if err := s.stateMachine.SetInitState(tablestate.TABLE_STATE_INIT); err != nil {
		return err
	}

	s.gameStartTime = time.Now()

	if err := s.initializeHandCard(); err != nil {
		return err
	}
	// 广播游戏开始消息
	msg := &protocol.ScGameStart{
		DiePoint: rand.Int31n(6) + 1,
		Games:    int32(s.games),
	}
	if err := s.BroadCast(protocol.PROTOID_SC_GAME_START, msg); err != nil {
		zlog.Errorf("broadCast game start failed, err=%s", err)
		return err
	}

	// 广播玩家手牌
	err := s.broadCastRaw(protocol.PROTOID_SC_CARD_INFO, func(pid uint64, seat int) proto.Message {
		return s.PackTableCardForPlayer(pid)
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
	ply := tableplayer.NewTablePlayer(plyData, s, win.NewGeneralWin(MAX_HAND_CARD_NUM))
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
	return !s.gameStartTime.IsZero() && s.gameStartTime.After(s.gameEndTime)
}

func (s *ScCardTable) IsGameEnd() bool {
	return !s.gameEndTime.IsZero() && s.gameEndTime.After(s.gameStartTime)
}

func (s *ScCardTable) IsEnd() bool {
	return !s.tableEndTime.IsZero()
}

func (s *ScCardTable) GetCreateTime() time.Time {
	return s.createTime
}

func (s *ScCardTable) GetTableNumber() uint32 {
	return s.id
}

// 广播同样的消息给所有玩家
func (s *ScCardTable) BroadCast(protoID protocol.PROTOID, msg proto.Message) error {
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
	reply.StartTime = s.gameStartTime.Unix()
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
	if err := s.BroadCast(protocol.PROTOID_SC_JOIN_TABLE, msg); err != nil {
		return err
	}

	return nil
}

const (
	HAND_CARD_NUM     = 13
	MAX_HAND_CARD_NUM = 14
)

func (s *ScCardTable) initializeHandCard() error {
	s.board = s.boardRule.NewBoard()
	s.shuffleRule.Shuffle(s.board.Cards)

	for i := 0; i < s.data.MaxPlayer; i++ {
		_ = s.players[i].InitHandCard(s.board.Cards[:HAND_CARD_NUM], MAX_HAND_CARD_NUM)
		s.board.Cards = s.board.Cards[HAND_CARD_NUM:]
	}
	return nil
}

func (s *ScCardTable) PackTableCardForPlayer(pid uint64) *protocol.ScCardInfo {
	msg := &protocol.ScCardInfo{}
	msg.TableCard = &protocol.TableCardData{}
	msg.TableCard.Total = int32(len(s.board.CardsTotal))
	msg.TableCard.Left = int32(len(s.board.Cards))

	for _, ply := range s.players {
		msg.CardInfo = append(msg.CardInfo, s.PackCardInfoForPlayer(ply, pid))
	}

	return msg
}

func (s *ScCardTable) PackCardInfoForPlayer(ply *tableplayer.TablePlayer, pid uint64) *protocol.PlayerCardInfo {
	var cards []int
	if ply.Pid == pid {
		cards = ply.Hcard.GetSortHandCard()
	} else {
		cards = ply.Hcard.GetGuardHandCard()
	}
	playerCardInfo := &protocol.PlayerCardInfo{
		HandCard: &protocol.HandCardData{},
	}

	for _, c := range cards {
		playerCardInfo.HandCard.Holders = append(playerCardInfo.HandCard.Holders, int32(c))
	}
	playerCardInfo.HandCard.Draw = int32(ply.Hcard.GetLastDraw())
	for _, k := range ply.Hcard.KongCards {
		playerCardInfo.HandCard.Kong = append(playerCardInfo.HandCard.Kong, &protocol.KongInfo{
			Card:     int32(k.Card),
			KongType: int32(k.KType),
		})
	}
	for _, p := range ply.Hcard.PongCards {
		playerCardInfo.HandCard.Pong = append(playerCardInfo.HandCard.Pong, int32(p))
	}
	for _, d := range ply.Hcard.DiscardCards {
		playerCardInfo.OutCard = append(playerCardInfo.OutCard, int32(d))
	}
	return playerCardInfo
}

func (s *ScCardTable) Update(delta time.Duration) {
	s.events.FireAll()
	if s.IsGameEnd() { // 游戏结束后不再处理timer和状态机
		return
	}

	// 更新玩家operate的timer
	s.UpdateOperateTimer() // 超时后的默认动作有可能会更新状态机的数据，所以要先做

	// 更新状态机
	if err := s.stateMachine.Update(); err != nil {
		zlog.Errorf("update state machine failed, err:%s", err)
	}

}

func (s *ScCardTable) DrawCard() error {
	turnPly := s.GetTurnPlayer()
	card, err := s.board.DrawForward()
	if err != nil {
		return err
	}
	zlog.Infof("player draw card, pid:%d, card:%d", turnPly.Pid, card)
	drawCmd := tableoperate.NewOperateDraw(card)
	if err = turnPly.DoOperate(drawCmd); err != nil {
		return err
	}

	err = s.broadCastRaw(protocol.PROTOID_SC_DRAW_CARD, func(pid uint64, seat int) proto.Message {
		msg := &protocol.ScDrawCard{Pid: turnPly.Pid, Card: -1}
		if turnPly.Pid == pid { // 摸到的牌只会发给本人
			msg.Card = int32(card)
		}
		return msg
	})
	if err != nil {
		return err
	}
	// 下发玩家操作
	if err = s.AfterPlyOperate(turnPly.Pid, drawCmd); err != nil {
		return err
	}
	return nil
}

func (s *ScCardTable) GetPlayers() []*tableplayer.TablePlayer {
	return s.players
}

func (s *ScCardTable) GetDiscardRule() irule.IDiscard {
	return s.discardRule
}

func (s *ScCardTable) GetTingRule() irule.ITing {
	return s.tingRule
}

func (s *ScCardTable) GetScoreModel() irule.IScoreCardModel {
	return s.scoreCardModel
}

func (s *ScCardTable) GetWinModeModel() irule.IWinModeModel {
	return s.winModeModel
}

func (s *ScCardTable) GetScorePointRule() irule.IScorePoint {
	return s.scorePointRule
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

	if err := ply.DoOperate(operate); err != nil {
		return err
	}

	if err := s.stateMachine.GetCurState().OnPlyOperate(pid, operate); err != nil {
		return err
	}

	// 更新下一个回合的玩家
	if operate.OpType != tableoperate.OPERATE_PASS && operate.OpType != tableoperate.OPERATE_DING_QUE {
		turnSeat := s.GetTurnSeat()
		maxPlayer := s.data.MaxPlayer
		plySeat := s.GetPlayerSeat(ply.Pid)
		// 靠后的玩家才会更新
		zlog.Debugf("try to update next turn seat, plyseat=%d, turnSeat=%d, maxPlayer=%d", plySeat, turnSeat, maxPlayer)
		if util.SeatRelative(plySeat, turnSeat, maxPlayer) > util.SeatRelative(s.nextTurnSeat, turnSeat, maxPlayer) {
			s.SetNextSeat(plySeat)
		}
	}

	if err := s.AfterPlyOperate(pid, operate); err != nil {
		return err
	}

	pbOperate := &protocol.ScNotifyOperate{
		Pid: pid,
		OpCmd: &protocol.OperateCommand{
			OpType: int32(operate.OpType),
			Data: &protocol.OperateData{
				Card: int32(operate.OpData.Card),
			},
		},
	}
	if err := s.BroadCast(protocol.PROTOID_SC_NOTIFY_OPERATE, pbOperate); err != nil {
		return err
	}

	cards := ply.Hcard.GetSortHandCard()
	zlog.Infof("after operate, card:%v, op:%v", cards, operate)

	// 检测本局是否结束
	s.CheckGameEnd()

	return nil
}

func (s *ScCardTable) CheckGameEnd() {
	var gamePlys []irule.GamePlayer
	for _, ply := range s.players {
		gamePlys = append(gamePlys, ply)
	}
	if s.gameModeRule.IsGameEnd(gamePlys, s.board) {
		s.OnGameEnd()
	}
}

func (s *ScCardTable) OnGameEnd() {
	s.gameEndTime = time.Now()
	zlog.Info("game end")

	// 广播通知结算
	msg := &protocol.ScGameEnd{
		Games: int32(s.games),
	}

	for _, ply := range s.players {
		summaryInfo := s.PackGameSummaryInfo(ply)
		msg.Summary = append(msg.Summary, summaryInfo)
	}
	_ = s.BroadCast(protocol.PROTOID_SC_GAME_END, msg)

	for _, ply := range s.players {
		ply.OnGameEnd()
	}
	// 清除庄家信息
	master := s.GetPlayerByPid(s.dealer)
	master.RemoveIdentity(gamedefine.TABLE_IDENTIY_DEALER)
	s.dealer = 0
	s.discards = []int{}

	if s.games >= int(s.data.GameTurn) {
		s.OnTableEnd()
		return
	}
}

func (s *ScCardTable) PackGameSummaryInfo(ply *tableplayer.TablePlayer) *protocol.ScGameSummaryInfo {
	summaryInfo := &protocol.ScGameSummaryInfo{
		WinInfo: []*protocol.ScWinInfo{},
	}
	summaryInfo.CardInfo = s.PackCardInfoForPlayer(ply, ply.Pid)
	for _, win := range ply.Wins {
		singleWin := &protocol.ScWinInfo{}
		for _, m := range win.Models {
			singleWin.CardModel = append(singleWin.CardModel, int32(m))
		}

		singleWin.WinMode = int32(win.Mode)

		singleWin.WinCard = int32(win.Card)
		summaryInfo.Score += int32(win.Score.GetFinalPoint())
		singleWin.Point = int32(win.Score.Point)
		singleWin.WinTime = win.Tm.UnixNano()
		summaryInfo.WinInfo = append(summaryInfo.WinInfo, singleWin)
	}

	for _, lose := range ply.Loses {
		summaryInfo.Score -= int32(lose.Score.GetFinalPoint())
	}
	return summaryInfo
}

func (s *ScCardTable) PackTableSummaryInfo() *protocol.ScTableEnd {
	summaryInfo := &protocol.ScTableEnd{}
	for _, ply := range s.players {
		gamesInfo := ply.GetGamesInfo()
		stat := &protocol.ScPlayerSummaryStatistics{}
		for _, game := range gamesInfo {
			// 统计胡牌和得分
			for _, w := range game.Wins {
				if winmode.IsDrawWin(w.Mode) {
					stat.DrawWin += 1
				} else {
					stat.DiscardWin += 1
				}
				stat.TotalScore += int32(w.Score.GetFinalPoint())
			}
			// 统计放炮和扣分
			for _, l := range game.Loses {
				stat.DiscardLose += 1
				stat.TotalScore -= int32(l.Score.GetFinalPoint())
			}
			// 统计杠
			for _, k := range game.Hcard.KongCards {
				switch k.KType {
				case gamedefine.KONG_TYPE_CONCEALED:
					stat.ConcealedKong += 1
				case gamedefine.KONG_TYPE_EXPOSED:
					stat.ExposedKong += 1
				case gamedefine.KONG_TYPE_RAIN:
					stat.RainKong += 1
				}
			}
			// todo 查叫
		}
		summaryInfo.Statistics = append(summaryInfo.Statistics, stat)
	}
	return summaryInfo
}

func (s *ScCardTable) OnTableEnd() {
	zlog.Info("table end")
	// 通知结算消息
	s.tableEndTime = time.Now()
	summaryInfo := s.PackTableSummaryInfo()
	_ = s.BroadCast(protocol.PROTOID_SC_TABLE_END, summaryInfo)
}

func (s *ScCardTable) AfterPlyOperate(pid uint64, operate tableoperate.OperateCommand) error {
	var err error
	switch operate.OpType {
	case tableoperate.OPERATE_DISCARD:
		err = s.AfterDiscard(pid, operate.OpData.Card)
	case tableoperate.OPERATE_KONG_CONCEALED: // 玩家可能会抢杠
		err = s.AfterConcealedKong(pid, operate.OpData.Card)
	case tableoperate.OPERATE_DISCARD_WIN:
		err = s.AfterDiscardWin(pid, operate.OpData.Card)
	case tableoperate.OPERATE_DRAW_WIN:
		err = s.AfterDrawWin(pid, operate.OpData.Card)
	case tableoperate.OPERATE_PONG:
		err = s.AfterPong(pid, operate.OpData.Card)
	case tableoperate.OPERATE_DRAW:
		err = s.AfterDraw(pid, operate.OpData.Card)
	}
	return err
}

func (s *ScCardTable) GetTimeoutOperate(pid uint64, ops []tableoperate.OperateCommand) tableoperate.OperateCommand {
	for _, op := range ops {
		switch op.OpType {
		case tableoperate.OPERATE_DISCARD: // 如果有出牌操作默认就是出牌
			return op
		case tableoperate.OPERATE_PASS: // 如果有跳过操作 那么默认就是跳过
			return op
		}
	}
	return tableoperate.NewOperatePass()
}

func (s *ScCardTable) UpdateOperateTimer() {
	// 每秒更新一次
	select {
	case <-s.secondTicker.C:
	default:
		return
	}
	for _, ply := range s.players {
		ops := ply.GetOperates()
		if ply.IsOperateTimeOut(s.GetOperateTimeout(ply.Pid, ops)) {
			zlog.Debugf("ply operate timeout, pid:%d ops:%v", ply.Pid, ops)
			defaultOp := s.GetTimeoutOperate(ply.Pid, ops)
			if err := s.OnPlyOperate(ply.Pid, defaultOp); err != nil {
				zlog.Errorf("do operate failed, err:%s", err)
			}
		}
	}
}

func (s *ScCardTable) GetOperateTimeout(pid uint64, ops []tableoperate.OperateCommand) time.Duration {
	//  return 10 * time.Second
	return -1
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

func (s *ScCardTable) UpdateTurnSeat() {
	s.turnSeat = s.nextTurnSeat
	updateMsg := &protocol.ScUpdateTurnPly{}
	updateMsg.Turn = int32(s.turnSeat)
	_ = s.BroadCast(protocol.PROTOID_SC_UPDATE_TURN_PLAYER, updateMsg)
}

func (s *ScCardTable) SetNextSeat(seat int) {
	if seat >= s.data.MaxPlayer {
		seat = 0
	}
	zlog.Debugf("set next TurnSeat from:%d->%d", s.nextTurnSeat, seat)
	s.nextTurnSeat = seat
}

func (s *ScCardTable) AfterDiscard(pid uint64, c int) error {
	s.discards = append(s.discards, c)
	for i := range s.players {
		ply := s.players[i]
		if ply.Pid == pid {
			continue
		}
		ops := ply.GetOperateWithDiscard(c)
		ply.SetOperate(ops)
		if err := s.distributeOperate(ply); err != nil {
			return err
		} // 下发玩家可以做的操作
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
		ply.SetOperate(ops)
		if err := s.distributeOperate(ply); err != nil {
			return err
		}

	}
	return nil
}

func (s *ScCardTable) AfterDrawWin(pid uint64, c int) error {
	ply := s.GetPlayerByPid(pid)
	winInfo := ply.GetLastWinInfo()
	if winInfo == nil {
		return errors.Errorf("can't find win info, pid:%d", pid)
	}
	// 自摸
	for _, ply := range s.players {
		if ply.Pid == pid {
			continue
		}
		ply.LoseByDiscard(winInfo.Score, pid, winInfo.Card, winInfo.Mode)
	}
	return nil
}

func (s *ScCardTable) AfterDiscardWin(pid uint64, c int) error {
	ply := s.GetPlayerByPid(pid)
	winInfo := ply.GetLastWinInfo()
	if winInfo == nil {
		return errors.Errorf("can't find win info, pid:%d", pid)
	}
	loser := s.GetPlayerByPid(winInfo.Loser)
	loser.LoseByDiscard(winInfo.Score, pid, winInfo.Card, winInfo.Mode)
	return nil
}

func (s *ScCardTable) AfterPong(pid uint64, c int) error {
	ply := s.GetPlayerByPid(pid)
	ply.SetOperate(ply.GetOperateWithPong(c))
	return s.distributeOperate(ply)
}

func (s *ScCardTable) AfterDraw(pid uint64, c int) error {
	ply := s.GetPlayerByPid(pid)
	ply.SetOperate(ply.GetOperateWithDraw(c))
	return s.distributeOperate(ply)
}

func (s *ScCardTable) DingQueFinish() error {

	// 通知定缺消息
	msg := &protocol.ScDingQueFinish{}
	for _, ply := range s.players {
		dqInfo := &protocol.ScDingqueInfo{
			Pid:    ply.Pid,
			DqSuit: int32(ply.Hcard.DingQueSuit),
		}
		msg.Dingque = append(msg.Dingque, dqInfo)
	}

	_ = s.BroadCast(protocol.PROTOID_SC_DING_QUE_FINISH, msg)

	turnPly := s.GetTurnPlayer()
	turnPly.SetOperate(turnPly.GetOperateWithDraw(turnPly.Hcard.GetLastDraw()))
	return s.distributeOperate(turnPly)
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
	_ = s.BroadCast(protocol.PROTOID_SC_PLAYER_READY, readData)

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

func (s *ScCardTable) GetWinMode(pid uint64, selfDraw bool) int {
	turnPly := s.GetTurnPlayer()
	info := irule.WinModeInfo{
		WinPid:   pid,
		DrawWin:  selfDraw,
		TurnOps:  turnPly.GetOperateLog(),
		TurnDraw: turnPly.Hcard.DrawCards,
		Dealer:   s.dealer,
		Discards: s.discards,
	}
	return s.winModeModel.GetWinMode(info)
}

func (s *ScCardTable) distributeOperate(ply *tableplayer.TablePlayer) error {
	ops := ply.GetOperates()
	if len(ops) == 0 {
		return nil
	}
	opdata := &protocol.ScPlayerOperate{}
	for _, op := range ops {
		opdata.OpCmd = append(opdata.OpCmd, &protocol.OperateCommand{
			OpType: int32(op.OpType),
			Data: &protocol.OperateData{
				Card: int32(op.OpData.Card),
			},
		})
	}
	if err := util.SendMsg(ply.Pid, protocol.PROTOID_SC_PLAYER_OPERATE, opdata); err != nil {
		return err
	}
	return nil
}

func (s *ScCardTable) GetBoardCard() *boardcard.BoardCard {
	return s.board
}

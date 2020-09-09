package sccardtable

import (
	"github.com/pkg/errors"
	"math/rand"
	"time"
	"zinx-mj/game/card/boardcard"
	"zinx-mj/game/gamedefine"
	"zinx-mj/game/rule/irule"
	"zinx-mj/game/table/tableplayer"
	"zinx-mj/mjerror"
	"zinx-mj/network/protocol"
	"zinx-mj/player"
	"zinx-mj/util"

	"github.com/aceld/zinx/zlog"

	"github.com/golang/protobuf/proto"
)

type ScCardTable struct {
	id      uint32                     // 桌子ID
	players []*tableplayer.TablePlayer // 房间的玩家
	startTm int64

	data           *ScTableData
	curPlayerIndex int
	event          *ScTableEvent
	board          *boardcard.BoardCard

	boardRule   irule.IBoard
	chowRule    irule.IChow
	discardRule irule.IDiscard
	drawRule    irule.IDraw
	kongRule    irule.IKong
	pongRule    irule.IPong
	shuffleRule irule.IShuffle
	tingRule    irule.ITing
	winRule     irule.IWin
	dealRule    irule.IDeal
}

func (s *ScCardTable) GetID() uint32 {
	return s.id
}

func (s *ScCardTable) Operate(operate irule.IOperate) error {
	panic("implement me")
}

func (s *ScCardTable) onGameStart() error {
	//handsCard := s.GetInitHandCard()
	// todo 下发手牌
	// 初始化庄家
	s.curPlayerIndex = rand.Intn(s.data.MaxPlayer)
	// 庄家摸牌
	dealer := s.players[s.curPlayerIndex]
	if err := s.Draw(dealer); err != nil {
		return err
	}
	return nil
}

func (s *ScCardTable) GetPlayer(pid player.PID) *tableplayer.TablePlayer {
	for _, ply := range s.players {
		if pid == ply.Pid {
			return ply
		}
	}
	return nil
}

func (s *ScCardTable) Join(plyData *tableplayer.TablePlayerData, identity uint32) (*tableplayer.TablePlayer, error) {
	ply := tableplayer.NewTablePlayer(plyData)
	ply.AddIdentity(identity)
	s.players = append(s.players, ply)

	s.event.Add(EVENT_JOIN, ply.Pid)

	return ply, nil
}

func (s *ScCardTable) IsFull() bool {
	return len(s.players) >= s.data.MaxPlayer
}

func (s *ScCardTable) Quit(pid player.PID) error {
	panic("implement me")
}

func (s *ScCardTable) GetStartTime() int64 {
	return s.startTm
}

func (s *ScCardTable) GetTableNumber() uint32 {
	return s.id
}

func (s *ScCardTable) BroadCast(protoID protocol.PROTOID, msg proto.Message) {
	for _, ply := range s.players {
		if err := util.SendMsg(ply.Pid, protoID, msg); err != nil {
			zlog.Errorf("braodcast to player failed, pid=%d, protoID=%d", ply.Pid, protoID)
		}
	}
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

func (s *ScCardTable) PackToPBMsg() proto.Message {
	reply := &protocol.ScScmjTableInfo{}
	reply.TableId = s.id
	reply.StartTime = s.GetStartTime()
	reply.Rule = s.data.PackToPBMsg().(*protocol.ScmjRule)
	for _, ply := range s.players {
		reply.Players = append(reply.Players, s.PackPlayerData(ply))
	}
	return reply
}

func NewTable(tableID uint32, master *tableplayer.TablePlayerData, data *ScTableData) (*ScCardTable, error) {
	t := &ScCardTable{
		id:      tableID,
		startTm: time.Now().Unix(),
		data:    data,
	}
	t.InitEvent()
	_, err := t.Join(master, gamedefine.TABLE_IDENTIY_MASTER|gamedefine.TABLE_IDENTIY_PLAYER)
	if err != nil {
		return t, err
	}
	return t, nil
}

func (s *ScCardTable) InitEvent() {
	s.event = NewScTableEvent()
	s.event.On(EVENT_JOIN, s.onJoinEvent)
	s.event.On(EVENT_GAME_START, s.onGameStart)
}

func (s *ScCardTable) onJoinEvent(pid player.PID) error {
	// 通知其它玩家该玩家加入了房间 // 其实应该延迟一帧发送，需要等待其他协议
	msg := &protocol.ScJoinTable{}
	ply := s.GetPlayer(pid)
	if ply == nil {
		return errors.WithStack(mjerror.ErrPlyNotFound)
	}
	msg.Player = s.PackPlayerData(ply)
	msg.SeatIndex = int32(len(s.players)) - 1
	s.BroadCast(protocol.PROTOID_SC_JOIN_TABLE, msg)

	// 人满了就开游戏
	if s.IsFull() {
		s.event.Add(EVENT_GAME_START)
	}
	return nil
}

func (s *ScCardTable) GetInitHandCard() [][]int {
	const HAND_CARD_NUM = 13
	s.board = s.boardRule.NewBoard()
	s.shuffleRule.Shuffle(s.board.Cards)
	var plyCards [][]int
	for i := 0; i < s.data.MaxPlayer; i++ {
		plyCards = append(plyCards, s.board.Cards[:HAND_CARD_NUM])
		s.board.Cards = s.board.Cards[HAND_CARD_NUM:]
	}
	return plyCards
}

func (s *ScCardTable) Update(delta time.Duration) {
	s.event.FireAll() // 处理所有的事件
}

func (s *ScCardTable) Draw(ply *tableplayer.TablePlayer) error {
	card, err := s.board.DrawForward()
	if err != nil {
		return errors.WithStack(err)
	}
	if err = s.drawRule.Draw(ply.PlyCard, card); err != nil {
		return errors.WithStack(err)
	}
	// todo 通知玩具摸牌消息
	return nil
}

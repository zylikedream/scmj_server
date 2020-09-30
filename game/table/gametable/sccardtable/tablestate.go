package sccardtable

import "zinx-mj/player"

const (
	STATE_EVENT_DRAW    = "state_event_draw"
	STATE_EVENT_DISCARD = "state_event_discard"
	STATE_EVENT_PASS    = "state_event_pass"
	STATE_EVENT_PONG    = "state_event_pong"
	STATE_EVENT_KONG    = "state_event_kong"
	STATE_EVENT_WIN     = "state_event_win"
)

const (
	TABLE_STATE_INIT      = "state_init"
	TABLE_STATE_DRAW      = "state_draw"
	TABLE_STATE_WAIT_WIN  = "state_discard"
	TABLE_STATE_WAIT_KONG = "state_discard"
	TABLE_STATE_WAIT_PONG = "state_discard"
)

func (s *ScCardTable) enterDrawCard() error {
	return s.drawCard(s.curPlayerIndex)
}

func (s *ScCardTable) enterWaitWin(pid player.PID, card int) error {
	for _, ply := range s.players {
		if ply.Pid == pid {
			continue
		}
		if ply.PlyCard.IsTingCard(card) {
			// todo 发送听牌操作
		}
	}
	return nil
}

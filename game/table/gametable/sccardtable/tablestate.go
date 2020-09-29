package sccardtable

import (
	"zinx-mj/player"

	"github.com/jinzhu/gorm"
)

const (
	STATE_EVENT_DRAW    = "state_event_draw"
	STATE_EVENT_DISCARD = "state_event_discard"
	STATE_EVENT_PASS    = "state_event_pass"
	STATE_EVENT_PONG    = "state_event_pong"
	STATE_EVENT_KONG    = "state_event_kong"
	STATE_EVENT_WIN     = "state_event_win"
	STATE_EVENT_EMPTY   = "state_event_empty"
)

const (
	TABLE_STATE_INIT      = "state_init"    // 摸牌
	TABLE_STATE_DRAW      = "state_draw"    // 摸牌
	TABLE_STATE_DISCARD   = "state_discard" // 出牌
	TABLE_STATE_WAIT_WIN  = "state_discard" // 等待胡牌
	TABLE_STATE_WAIT_KONG = "state_discard" // 等待杠牌
	TABLE_STATE_WAIT_PONG = "state_discard" // 等待碰牌
)

type ParamWithCardAndTable struct {
	t    *ScCardTable
	card int
	pid  player.PID
}

func OnDrawEnter(t interface{}, tx *gorm.DB) error {
	tb := t.(*ScCardTable)
	return tb.drawCard(tb.curPlayerIndex)
}

func OnDiscardEnter(t interface{}, tx *gorm.DB) error {
	return nil
}

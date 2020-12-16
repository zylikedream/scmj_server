package sccardtable

import (
	"sync"

	"github.com/aceld/zinx/zlog"
	"github.com/sadlil/go-trigger"
)

const (
	EVENT_JOIN       = "event_join"
	EVENT_GAME_START = "event_game_start"
)

type event struct {
	name   string
	params []interface{}
}

type ScTableEvent struct {
	trigger.Trigger
	events []event
	mu     sync.Mutex
}

func NewScTableEvent() *ScTableEvent {
	return &ScTableEvent{
		Trigger: trigger.New(),
	}
}

func (s *ScTableEvent) Add(name string, params ...interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, event{name, params})
}

func (s *ScTableEvent) FireAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, eve := range s.events {
		_, err := s.FireBackground(eve.name, eve.params...)
		if err != nil {
			zlog.Errorf("fire error:%s", err)
		}
	}
	s.events = []event{}
}

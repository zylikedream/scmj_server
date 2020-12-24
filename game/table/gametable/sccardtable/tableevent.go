package sccardtable

import (
	"reflect"

	"github.com/aceld/zinx/zlog"
	"github.com/pkg/errors"
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
	events   []event
	callback map[string]interface{}
}

func NewScTableEvent() *ScTableEvent {
	return &ScTableEvent{}
}

func (s *ScTableEvent) Add(name string, params ...interface{}) {
	s.events = append(s.events, event{name, params})
}

func (s *ScTableEvent) Register(name string, callback interface{}) error {
	if _, ok := s.callback[name]; ok {
		return errors.Errorf("event already defined")
	}
	if reflect.ValueOf(callback).Type().Kind() != reflect.Func {
		return errors.Errorf("callback is not a function")
	}
	s.callback[name] = callback
	return nil
}

func (s *ScTableEvent) FireAll() {
	for _, eve := range s.events {
		f, in, err := s.parse(eve.name, eve.params...)
		if err == nil {
			zlog.Errorf("fire event failed, event:%v", eve)
		}
		f.Call(in)
	}
	s.events = []event{}
}

func (s *ScTableEvent) parse(event string, params ...interface{}) (reflect.Value, []reflect.Value, error) {
	cb, ok := s.callback[event]
	if !ok {
		return reflect.Value{}, nil, errors.Errorf("no callback found for event")
	}
	f := reflect.ValueOf(cb)
	if len(params) != f.Type().NumIn() {
		return reflect.Value{}, nil, errors.Errorf("parameter mismatched")
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f, in, nil
}

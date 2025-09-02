package eventify

import (
	"encoding/json"
	"reflect"
)

// Namable is an interface that can be used to name a listener
// Only listeners that implement this interface will be unregistered
type Namable interface {
	Name() string
}

type Listener interface {
	Handle(event Event) error
}

func NewListener(handle func(event Event) error) Listener {
	if handle == nil {
		handle = func(Event) error { return nil }
	}
	return &listener{
		handle: handle,
	}
}

type listener struct {
	handle func(event Event) error
}

func (l *listener) Handle(event Event) error {
	return l.handle(event)
}

func (l *listener) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"handle": reflect.TypeOf(l.handle).String(),
	})
}

func NewNamedListener(name string, handle func(event Event) error) Listener {
	if handle == nil {
		handle = func(Event) error { return nil }
	}
	return &namedListener{
		name:   name,
		handle: handle,
	}
}

type namedListener struct {
	name   string
	handle func(event Event) error
}

func (l *namedListener) Name() string {
	return l.name
}

func (l *namedListener) Handle(event Event) error {
	return l.handle(event)
}

func (l *namedListener) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"name":   l.name,
		"handle": reflect.TypeOf(l.handle).String(),
	})
}

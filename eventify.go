package eventify

import (
	"encoding/json"
	"sync"
)

// Eventify is a struct that represents an event emitter.
type Eventify struct {
	listeners sync.Map
	mutex     sync.RWMutex
	log       Log
}

// New creates a new Eventify instance with the default logger.
func New() *Eventify {
	return NewEventify()
}

// NewEventifyWithLog creates and returns a new Eventify instance with the provided logger.
func NewEventifyWithLog(log Log) *Eventify {
	return NewEventify(WithLogger(log))
}

// NewEventify creates and returns a new Eventify instance with the provided logger.
// The logger is used for debugging and error reporting throughout the Eventify instance's lifecycle.
// If no logger is needed, you can pass nil, but it's recommended to provide a logger for better observability.
func NewEventify(opts ...OptionFunc) *Eventify {
	o := NewOption(opts...)
	ev := &Eventify{
		listeners: sync.Map{},
		mutex:     sync.RWMutex{},
		log:       o.log,
	}
	return ev
}

// Register adds an event listener for the specified event type.
// The listener will be called whenever an event of the matching type is emitted.
// Multiple listeners can be registered for the same event type.
// This method is thread-safe.
func (e *Eventify) Register(eventType string, listener Listener) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	listeners, _ := e.listeners.LoadOrStore(eventType, []Listener{})
	e.listeners.Store(eventType, append(listeners.([]Listener), listener))
	e.log.Debug("eventify register", "event_type", eventType, "listener", listener)
}

// Unregister removes event listeners for the specified event type.
// If no listeners are provided, all listeners for the event type are removed.
// If specific listeners are provided, only those listeners will be removed.
// This method is thread-safe.
func (e *Eventify) Unregister(eventType string, listeners ...Listener) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if len(listeners) == 0 {
		e.listeners.Delete(eventType)
		return
	}
	e.log.Debug("eventify unregister", "event_type", eventType, "listeners", listeners)

	namedListeners := []Namable{}
	for _, listener := range listeners {
		if namable, ok := listener.(Namable); ok {
			namedListeners = append(namedListeners, namable)
		}
	}
	if len(namedListeners) == 0 {
		return
	}
	ls, ok := e.listeners.Load(eventType)
	if !ok {
		return
	}
	lsCopy := []Namable{}
	newLs := []Listener{}
	for _, listener := range ls.([]Listener) {
		if namable, ok := listener.(Namable); ok {
			lsCopy = append(lsCopy, namable)
		} else {
			newLs = append(newLs, listener)
		}
	}
	if len(lsCopy) == 0 {
		return
	}
	for _, listener := range lsCopy {
		for _, l := range namedListeners {
			if listener.Name() != l.Name() {
				newLs = append(newLs, listener.(Listener))
			}
		}
	}
	e.listeners.Store(eventType, newLs)
}

// Emit dispatches an event to all registered listeners for the event's type.
// The event is processed synchronously unless the event or listener implements IsAsync.
// If the event implements ErrorHandler, any errors from listeners will be handled asynchronously.
func (e *Eventify) Emit(event Event) {
	e._Emit(event)
}

// EmitBy creates and emits a new event with the specified type and payload.
// If the payload is already an Event, it will be emitted directly.
// Otherwise, a new event is created with the given type and payload.
// The payload will be automatically converted to bytes using JSON marshaling if needed.
func (e *Eventify) EmitBy(eventType string, payload any) {
	if event, ok := payload.(Event); ok {
		e._Emit(event)
		return
	}
	e._Emit(NewEvent(eventType, e._AnyToBytes(payload)))
}

func (e *Eventify) _Emit(event Event) {
	listeners := e._MatchedListeners(event.Type())
	_, isAsyncEvent := event.(IsAsync)
	for _, listener := range listeners {
		_, isAsyncListener := listener.(IsAsync)
		e._Trigger(event, listener, isAsyncEvent || isAsyncListener)
	}
	e.log.Debug("eventify emited", "event", event.Type(), "listeners", listeners)
}

func (e *Eventify) _Trigger(event Event, listener Listener, async bool) {
	errHandler, hasErrorHandler := event.(ErrorHandler)
	if async {
		go func() {
			if err := listener.Handle(event); err != nil && hasErrorHandler {
				go errHandler.ErrorHandler(event, err)
			}
		}()
		return
	}
	if err := listener.Handle(event); err != nil && hasErrorHandler {
		errHandler.ErrorHandler(event, err)
	}
}

func (e *Eventify) _MatchedListeners(eventType string) []Listener {
	listeners := make([]Listener, 0)
	e.listeners.Range(func(key, value any) bool {
		if NewMatcher(key.(string)).Match(eventType) {
			listeners = append(listeners, value.([]Listener)...)
		}
		return true
	})
	return listeners
}

func (e *Eventify) _AnyToBytes(payload any) []byte {
	if payload == nil {
		return nil
	}
	switch p := payload.(type) {
	case string:
		return []byte(p)
	case []byte:
		return p
	default:
		bz, err := json.Marshal(payload)
		if err != nil {
			return nil
		}
		return bz
	}
}

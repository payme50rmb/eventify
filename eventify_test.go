package eventify

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventify_RegisterAndEmit(t *testing.T) {
	tests := []struct {
		name        string
		eventType   string
		shouldPanic bool
		setup       func(*Eventify) (bool, func() bool)
	}{
		{
			name:      "basic event emission",
			eventType: "test.event",
			setup: func(e *Eventify) (bool, func() bool) {
				var triggered bool
				e.Register("test.event", NewListener(func(event Event) error {
					triggered = true
					assert.Equal(t, "test.event", event.Type())
					return nil
				}))
				return triggered, func() bool { return triggered }
			},
		},
		{
			name:      "multiple listeners",
			eventType: "multi.event",
			setup: func(e *Eventify) (bool, func() bool) {
				var count int
				for i := 0; i < 3; i++ {
					e.Register("multi.event", NewListener(func(event Event) error {
						count++
						return nil
					}))
				}
				return false, func() bool { return count == 3 }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEventify(nil)
			_, check := tt.setup(e)

			e._Emit(NewEvent(tt.eventType, nil))

			assert.True(t, check(), "event handlers not triggered correctly")
		})
	}
}

func TestEventify_ConcurrentAccess(t *testing.T) {
	e := NewEventify(nil)
	var wg sync.WaitGroup
	const numListeners = 100
	const numEmitters = 10

	// Register listeners concurrently
	for i := 0; i < numListeners; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.Register("concurrent.event", NewNamedListener("listener_"+string(rune(i)), nil))
		}(i)
	}

	// Emit events concurrently
	for i := 0; i < numEmitters; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.EmitBy("concurrent.event", []byte("test"))
		}()
	}

	// Wait for all goroutines to complete
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Add timeout to prevent hanging
	select {
	case <-done:
		// All goroutines completed
	case <-time.After(5 * time.Second):
		t.Fatal("test timed out, possible deadlock")
	}

	// Verify listeners were registered
	listeners := loadAllListeners(e)
	assert.Len(t, listeners["concurrent.event"], numListeners, "not all listeners were registered")
}

func loadAllListeners(e *Eventify) map[string][]Listener {
	listeners := map[string][]Listener{}
	e.listeners.Range(func(key, value any) bool {
		listeners[key.(string)] = value.([]Listener)
		return true
	})
	return listeners
}

func TestEventify_Unregister(t *testing.T) {
	type testCase struct {
		name     string
		setup    func(*Eventify) (eventType string, listeners []Listener)
		validate func(*testing.T, *Eventify, string, []Listener)
	}

	tests := []testCase{
		{
			name: "unregister all listeners for event",
			setup: func(e *Eventify) (string, []Listener) {
				eventType := "test.event"
				e.Register(eventType, NewListener(nil))
				e.Register(eventType, NewListener(nil))
				return eventType, nil
			},
			validate: func(t *testing.T, e *Eventify, eventType string, _ []Listener) {
				listeners := loadAllListeners(e)
				assert.Empty(t, listeners[eventType], "all listeners should be removed")
			},
		},
		{
			name: "unregister specific named listener",
			setup: func(e *Eventify) (string, []Listener) {
				eventType := "test.event"
				listener1 := NewNamedListener("l1", nil)
				listener2 := NewNamedListener("l2", nil)
				e.Register(eventType, listener1)
				e.Register(eventType, listener2)
				return eventType, []Listener{listener1}
			},
			validate: func(t *testing.T, e *Eventify, eventType string, toRemove []Listener) {
				listeners := loadAllListeners(e)
				require.Len(t, listeners[eventType], 1, "only one listener should remain")
				assert.Equal(t, "l2", listeners[eventType][0].(Namable).Name())
			},
		},
		{
			name: "unregister non-existent listener",
			setup: func(e *Eventify) (string, []Listener) {
				eventType := "test.event"
				e.Register(eventType, NewNamedListener("l1", nil))
				return eventType, []Listener{NewNamedListener("non-existent", nil)}
			},
			validate: func(t *testing.T, e *Eventify, eventType string, _ []Listener) {
				listeners := loadAllListeners(e)
				assert.Len(t, listeners[eventType], 1, "existing listener should remain")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEventify(nil)
			eventType, toRemove := tt.setup(e)

			e.Unregister(eventType, toRemove...)

			tt.validate(t, e, eventType, toRemove)
		})
	}
}

func TestEventify_EmitBy(t *testing.T) {
	t.Run("emits event with payload", func(t *testing.T) {
		e := NewEventify(nil)
		var receivedPayload []byte

		e.Register("test.event", NewListener(func(event Event) error {
			receivedPayload = event.Payload()
			return nil
		}))

		payload := []byte("test payload")
		e.EmitBy("test.event", payload)

		assert.Equal(t, payload, receivedPayload)
	})

	t.Run("handles nil payload", func(t *testing.T) {
		e := NewEventify(nil)
		var called bool

		e.Register("test.event", NewListener(func(event Event) error {
			called = true
			assert.Nil(t, event.Payload())
			return nil
		}))

		e.EmitBy("test.event", nil)

		assert.True(t, called, "listener should be called")
	})
}

func TestEventify_ErrorHandling(t *testing.T) {
	t.Run("error from listener is handled", func(t *testing.T) {
		e := NewEventify(nil)
		errChan := make(chan error, 1)

		e.Register("error.event", NewListener(func(event Event) error {
			return assert.AnError
		}))

		e.Register("error.event", NewListener(nil)) // This one succeeds

		e._Emit(&mockErrorEvent{errChan: errChan})

		select {
		case err := <-errChan:
			assert.Equal(t, assert.AnError, err)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("error handler not called")
		}
	})
}

type mockErrorEvent struct {
	errChan chan error
}

func (m *mockErrorEvent) Type() string    { return "error.event" }
func (m *mockErrorEvent) Payload() []byte { return nil }
func (m *mockErrorEvent) ErrorHandler(_ Event, err error) {
	m.errChan <- err
}

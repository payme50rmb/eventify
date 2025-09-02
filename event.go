package eventify

// ErrorHandler is an interface that can be implemented by events to handle errors that occur during event processing.
type ErrorHandler interface {
	ErrorHandler(event Event, err error)
}

// Event is an interface that represents an event.
type Event interface {
	Type() string
	Payload() []byte
}

// NewEvent creates a new event with the specified type and payload.
func NewEvent(eventType string, payload []byte) Event {
	return &event{
		eventType: eventType,
		payload:   payload,
	}
}

type event struct {
	eventType string
	payload   []byte
}

func (e *event) Type() string {
	return e.eventType
}

func (e *event) Payload() []byte {
	return e.payload
}

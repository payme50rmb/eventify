package eventify

type ErrorHandler interface {
	ErrorHandler(event Event, err error)
}

type Event interface {
	Type() string
	Payload() []byte
}

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

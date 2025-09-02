package eventify

// Namable is an interface that can be used to name a listener
// Only listeners that implement this interface will be unregistered
type Namable interface {
	Name() string
}

// Listener is an interface that represents an event listener.
type Listener interface {
	Handle(event Event) error
}

// NewListener creates a new listener with the specified handle function.
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

// NewNamedListener creates a new named listener with the specified name and handle function.
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

package eventify

// IsAsync is an interface that can be implemented by events and listeners to indicate that they should be processed asynchronously.
type IsAsync interface {
	isAsync()
}

// IAmAsync is a struct that implements the IsAsync interface.
type IAmAsync struct {
}

func (IAmAsync) isAsync() {}

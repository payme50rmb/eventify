package eventify

type IsAsync interface {
	isAsync()
}

type IAmAsync struct {
}

func (IAmAsync) isAsync() {}

package eventify

type Log interface {
	Debug(msg string, kvs ...any)
}

type NoLog struct{}

func (*NoLog) Debug(msg string, kvs ...any) {}

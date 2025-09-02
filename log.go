package eventify

// Log is an interface that represents a logger.
type Log interface {
	Debug(msg string, kvs ...any)
}

// NoLog is a logger that does nothing.
type NoLog struct{}

// Debug does nothing.
func (*NoLog) Debug(msg string, kvs ...any) {}

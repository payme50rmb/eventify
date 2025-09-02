package eventify

// Option is a struct that represents an option for the Eventify instance.
type Option struct {
	log Log
}

// OptionFunc is a function that configures an Option.
type OptionFunc func(*Option)

// WithLogger sets the logger for the Eventify instance.
func WithLogger(log Log) OptionFunc {
	return func(o *Option) {
		o.log = log
	}
}

// NewOption creates a new Option with the specified options.
func NewOption(opts ...OptionFunc) *Option {
	o := &Option{
		log: &NoLog{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

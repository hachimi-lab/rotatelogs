package rotatelogs

import (
	"time"
)

const (
	DefaultTimePeriod = Daily
	DefaultMaxAge     = time.Hour * 24 * 7
)

var defaultOpts = options{
	timePeriod: DefaultTimePeriod,
	maxAge:     DefaultMaxAge,
}

type (
	Option interface {
		apply(*options)
	}
	options struct {
		timePeriod TimePeriod
		maxAge     time.Duration
	}
	optTimePeriod TimePeriod
	optMaxAge     time.Duration
)

func (opt optTimePeriod) apply(opts *options) {
	opts.timePeriod = TimePeriod(opt)
}

func (opt optMaxAge) apply(opts *options) {
	opts.maxAge = time.Duration(opt)
}

func WithTimePeriod(timePeriod TimePeriod) Option {
	if !timePeriod.isValid() {
		timePeriod = DefaultTimePeriod
	}
	return optTimePeriod(timePeriod)
}

func WithMaxAge(maxAge time.Duration) Option {
	if maxAge < time.Minute {
		maxAge = time.Minute
	}
	return optMaxAge(maxAge)
}

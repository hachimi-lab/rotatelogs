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
	Option  func(opts *options)
	options struct {
		timePeriod TimePeriod
		maxAge     time.Duration
	}
)

func WithTimePeriod(timePeriod TimePeriod) Option {
	if !timePeriod.isValid() {
		timePeriod = DefaultTimePeriod
	}
	return func(opts *options) {
		opts.timePeriod = timePeriod
	}
}

func WithMaxAge(maxAge time.Duration) Option {
	if maxAge < time.Minute {
		maxAge = time.Minute
	}
	return func(opts *options) {
		opts.maxAge = maxAge
	}
}

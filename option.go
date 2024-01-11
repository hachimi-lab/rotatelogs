package rotatelogs

import (
	"time"
)

const (
	DefaultFileRotateTime = EveryDay
	DefaultFileAge        = time.Hour * 24 * 7
)

var defaultOpts = options{
	rotateTime: DefaultFileRotateTime,
	maxAge:     DefaultFileAge,
}

type (
	Option interface {
		apply(*options)
	}
	options struct {
		rotateTime RotateTime
		maxAge     time.Duration
	}
	optRotateTime RotateTime
	optMaxAge     time.Duration
)

func (opt optRotateTime) apply(opts *options) {
	opts.rotateTime = RotateTime(opt)
}

func (opt optMaxAge) apply(opts *options) {
	opts.maxAge = time.Duration(opt)
}

func WithRotateTime(rotateTime RotateTime) Option {
	if rotateTime <= EveryMinute {
		rotateTime = EveryMinute
	} else if rotateTime >= EveryDay {
		rotateTime = EveryDay
	} else {
		rotateTime = EveryHour
	}
	return optRotateTime(rotateTime)
}

func WithMaxAge(maxAge time.Duration) Option {
	if maxAge < time.Minute {
		maxAge = time.Minute
	}
	return optMaxAge(maxAge)
}

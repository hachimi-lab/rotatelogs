package rotatelogs

import (
	"time"
)

const (
	DefaultFileRotateTime = time.Hour * 24
	DefaultFileAge        = time.Hour * 24 * 7
)

type options struct {
	rotateTime time.Duration
	maxAge     time.Duration
}

type Option interface {
	apply(*options)
}

type rotateTimeOpt time.Duration

func (opt rotateTimeOpt) apply(opts *options) {
	opts.rotateTime = time.Duration(opt)
}

type maxAgeOpt time.Duration

func (opt maxAgeOpt) apply(opts *options) {
	opts.maxAge = time.Duration(opt)
}

var defaultOpts = options{
	rotateTime: DefaultFileRotateTime,
	maxAge:     DefaultFileAge,
}

func WithRotateTime(rotateTime time.Duration) Option {
	if rotateTime < time.Minute {
		rotateTime = time.Minute
	}
	return rotateTimeOpt(rotateTime)
}

func WithMaxAge(maxAge time.Duration) Option {
	if maxAge < time.Minute {
		maxAge = time.Minute
	}
	return maxAgeOpt(maxAge)
}

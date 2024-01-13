package rotatelogs

import (
	"time"
)

const (
	DefaultFileRotateType = EveryDay
	DefaultFileAge        = time.Hour * 24 * 7
)

var defaultOpts = options{
	rotateType: DefaultFileRotateType,
	maxAge:     DefaultFileAge,
}

type (
	Option interface {
		apply(*options)
	}
	options struct {
		rotateType RotateType
		maxAge     time.Duration
	}
	optRotateType RotateType
	optMaxAge     time.Duration
)

func (opt optRotateType) apply(opts *options) {
	opts.rotateType = RotateType(opt)
}

func (opt optMaxAge) apply(opts *options) {
	opts.maxAge = time.Duration(opt)
}

func WithRotateType(rotateType RotateType) Option {
	if !rotateType.isValid() {
		rotateType = DefaultFileRotateType
	}
	return optRotateType(rotateType)
}

func WithMaxAge(maxAge time.Duration) Option {
	if maxAge < time.Minute {
		maxAge = time.Minute
	}
	return optMaxAge(maxAge)
}

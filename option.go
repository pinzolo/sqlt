package sqlt

import "time"

// Option is function to set optional behavior.
type Option func(*config)

var (
	// TimeFunc is option for setting custom time func.
	TimeFunc = func(fn func() time.Time) Option {
		return func(conf *config) {
			conf.timeFn = fn
		}
	}
)

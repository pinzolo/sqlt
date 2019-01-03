package sqlt

import "time"

// Option is function to set optional behavior.
type Option func(*config)

var (
	// TimeFunc is option for setting custom time func.
	TimeFunc = func(fn func() time.Time) Option {
		return func(conf *config) {
			conf.timeFunc = fn
		}
	}

	// Annotation is option for enabling annotative in template.
	Annotation = func() Option {
		return func(conf *config) {
			conf.annotative = true
		}
	}
)

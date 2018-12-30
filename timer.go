package sqlt

import "time"

type timer struct {
	fn         func() time.Time
	cache      time.Time
	cacheIndex int
	nowCnt     int
}

func newTimer(fn func() time.Time) *timer {
	if fn == nil {
		return &timer{
			fn:     time.Now,
			nowCnt: 1,
		}
	}
	return &timer{fn: fn, nowCnt: 1}
}

func (t *timer) cached() bool {
	return !t.cache.IsZero()
}

func (t *timer) time() time.Time {
	if !t.cached() {
		t.cache = t.fn()
	}
	return t.cache
}

func (t *timer) now() time.Time {
	t.nowCnt++
	return t.fn()
}

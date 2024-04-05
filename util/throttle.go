package util

import (
	"sync"
	"time"
)

type Throttler interface {
	Do(f func())
}

func NewThrottle(duration time.Duration) Throttler {
	return &throttle{
		duration: duration,
		once:     &sync.Once{},
	}
}

type throttle struct {
	duration time.Duration
	once     *sync.Once
	mutex    sync.Mutex
}

func (t *throttle) Do(f func()) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.once.Do(func() {
		// 有效时间内，仅执行一次调用
		time.Sleep(t.duration)
		f()
		// 为下次执行做好准备
		t.once = &sync.Once{}
	})
}

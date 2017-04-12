package utils

import (
	"time"
)

type Timer struct {
	*time.Timer
	dur time.Duration
}

func NewTimer(second int64) *Timer {
	t := &Timer{}
	t.dur = time.Duration(second)
	t.Timer = time.NewTimer(t.dur * time.Second)
	return t
}

func (t *Timer) Update(now uint64) bool {
	select {
	case <-t.C:
		t.Reset(t.dur * time.Second)
		return true
	default:
		return false
	}
}

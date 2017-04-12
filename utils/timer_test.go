package utils

import (
	"fmt"
	"testing"
	"time"
)

func Test_timer(t *testing.T) {
	timer_1 := NewTimer(1)
	f := uint64(time.Now().Unix())
	for {
		now := uint64(time.Now().Unix())
		if timer_1.Update(now) {
			fmt.Println("exec second %d", now)
		}
		if f+100 < now {
			break
		}
	}
}

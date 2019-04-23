package util

import "runtime"

// GetMemState 获取进程当前的内存状态
func GetMemState() (uint64, uint64) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return m.Sys, m.PauseTotalNs
}

package util

import "runtime"

func GetMemState() (uint64, uint64) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	return m.Sys, m.PauseTotalNs
}
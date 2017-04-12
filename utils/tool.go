package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenUUID() uint64 {
	rand.Seed(time.Now().UnixNano())
	return uint64(rand.Int63())
}

func GetDayBegin() time.Time {
	return time.Now().AddDate(0, 0, 1).Round(time.Hour * 24)
}

func ChangeStringToUint32(s string) uint32 {
	n, e := strconv.ParseUint(s, 10, 32)
	if e != nil {
		return 0
	}
	return uint32(n)
}

func ChangeUint32ToString(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}
func ChangeUint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}
func ChangeStringToArrayUint32(s string, sep string) []uint32 {
	if s == "" || s == "0" {
		return nil
	}

	var d []uint32
	tmp := strings.Split(s, sep)
	for _, v := range tmp {
		d = append(d, ChangeStringToUint32(v))
	}

	return d
}
func ChangeStringToMapUint32(s, sep1, sep2 string) map[uint32]uint32 {
	if s == "" || s == "0" {
		return nil
	}

	d := map[uint32]uint32{}
	tmp := strings.Split(s, sep1)
	for _, v := range tmp {
		tmp1 := ChangeStringToArrayUint32(v, sep2)
		if len(tmp1) != 2 {
			continue
		}
		d[tmp1[0]] = tmp1[1]
	}

	return d
}

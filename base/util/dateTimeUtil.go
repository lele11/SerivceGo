package util

import (
	"strings"
	"time"

	log "github.com/cihub/seelog"
)

//返回今天零点对应的时间戳
func GetTodayBeginStamp() int64 {
	t := time.Now()
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return t.Unix()
}

// 返回本周一零点对应的时间戳
func GetThisWeekBeginStamp() int64 {
	//TODO 思路复杂 只需要通过周内的 天数差 计算，不需要通过周数计算
	t := time.Now()
	year, week := t.ISOWeek()
	date := time.Date(year, 0, 0, 0, 0, 0, 0, time.Local)
	isoYear, isoWeek := date.ISOWeek()

	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	for isoYear < year {
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}

	for isoWeek < week {
		date = date.AddDate(0, 0, 1)
		_, isoWeek = date.ISOWeek()
	}

	return date.Unix()
}

//判断一个时间戳是周几(1-7)
func WhatDayOfWeek(stamp int64) int {
	weekDays := map[string]int{
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
		"sunday":    7,
	}

	tm := time.Unix(stamp, 0)
	day := strings.ToLower(tm.Weekday().String())

	return weekDays[day]
}

//获取本周某一日零点对应的时间戳，whatDay表示周几(1-7)
func GetWeekDayBeginStamp(whatDay uint8) int64 {
	begin := GetThisWeekBeginStamp()
	day := int64(24 * 60 * 60)
	return begin + int64(whatDay-1)*day
}

//判断一个时间戳是否在本周的时间段内
func IsInThisWeek(stamp int64) bool {
	begin := GetThisWeekBeginStamp()
	if stamp < begin {
		return false
	}

	day := int64(24 * 60 * 60)
	if stamp >= begin+7*day {
		return false
	}

	return true
}

//获取下一个整点的时间戳
func GetNextHoursStamp() int64 {
	now := time.Now()
	return now.Unix() + int64((60-now.Minute())*60-now.Second())
}

//获取下一个整点的时间差
func GetNextHoursDur() int64 {
	now := time.Now()
	return int64((60-now.Minute())*60 - now.Second())
}

func TimeStringToTime(ts string) *time.Time {
	t, e := time.ParseInLocation("2006|01|02|15|04|05", ts, time.Local)
	if e != nil {
		log.Error("ParseInLocation err: ", e, ts)
	}
	return &t
}
func TimeIsDayBegin() bool {
	if time.Now().Hour() == 0 &&
		time.Now().Minute() == 0 &&
		time.Now().Second() == 0 {
		return true
	}
	return false
}

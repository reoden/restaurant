package common

import (
	"time"
)

// 时间格式
const (
	TimeFormat1 = "2006-01-02 15:04:05"
	TimeFormat2 = "20060102150405"
	// ES @timestamp 常用格式
	TimeFormat3 = "2006-01-02T15:04:05.000Z"
	DateFormat1 = "2006-01-02"
	DateFormat2 = "20060102"
)

// NowDateStr 取得当天日期字符串，具体格式按layout获取
func NowDateStr(layout string) string {
	now := time.Now()
	return now.Format(layout)
}

// Diff 取得两个给定字符串的日期之差
func Diff(layout string, afterTimeStr string, beforeTimeStr string) time.Duration {
	afterTime, _ := time.Parse(layout, afterTimeStr)
	beforeTime, _ := time.Parse(layout, beforeTimeStr)
	return afterTime.Sub(beforeTime)
}

// MicroTs 取得13位的毫秒TS
func MicroTs() int {
	return int(time.Now().UnixNano() / 1000000)
}

// Ts 取得秒级TS
func Ts() int {
	return int(time.Now().Unix())
}

// ParseTimeLocal 将字符串转为time类
func ParseTimeLocal(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

// StrToTime 字符串转time类
func StrToTime(str string, format string) (time.Time, error) {
	if format == "" {
		format = TimeFormat1
	}
	time, err := time.ParseInLocation(format, str, time.Local)
	return time, err
}

func TimeToDateStr(t time.Time) string {
	return t.Format(DateFormat1)
}

func TimeToDatetimeStr(t time.Time) string {
	return t.Format(TimeFormat1)
}

// TSToTime 时间戳转time类
func TSToTime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

// TimeToStr 时间转字符串
func TimeToStr(t time.Time, format string) string {
	if format == "" {
		format = TimeFormat1
	}
	timeStr := t.Format(format)
	return timeStr
}

func NowInLocal() time.Time {
	local, _ := time.LoadLocation("Local")
	return time.Now().In(local)
}

func NowLocalStr(format string) string {
	return TimeToStr(NowInLocal(), format)
}

func GetMonday() time.Time {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	monday := now.AddDate(0, 0, offset)
	return monday
}

func EndOfDay(t time.Time) time.Duration {
	year, month, day := t.Date()
	endOfDay := time.Date(year, month, day, 23, 59, 59, 0, t.Location())
	return endOfDay.Sub(t)
}

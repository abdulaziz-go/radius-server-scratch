package timeUtil

import (
	"fmt"
	"time"
)

const DayTimestamp = 86400 // 24 * 60 * 60

func NowUnixTime() int64 {
	return time.Now().Unix()
}

func NowUnixMilliTime() int64 {
	return time.Now().UnixMilli()
}

func UnixMilliToUnix(timestamp int64) int64 {
	return timestamp / 1000
}

func Int64ToTime(value int64) time.Time {
	return time.Unix(value, 0)
}

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func SleepMs(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func DurationSeconds(value int) time.Duration {
	return time.Duration(value) * time.Second
}

func DurationMillisecond(value int) time.Duration {
	return time.Duration(value) * time.Millisecond
}

func IsUnixTime(timestamp int64) bool {
	t := time.Unix(timestamp, 0)
	return !t.IsZero()
}

func StartOfDayUnixUTC(timestamp int64) int64 {
	timestampUnix := time.Unix(timestamp, 0).UTC()
	startOfDay := time.Date(timestampUnix.Year(), timestampUnix.Month(), timestampUnix.Day(), 0, 0, 0, 0, time.UTC)
	return startOfDay.Unix()
}

func ConvertMillisecondsToTimeString(ms int64) string {
	seconds := float64(ms) / 1000
	hours := int(seconds) / 3600
	seconds = seconds - float64(hours*3600)
	minutes := int(seconds) / 60
	seconds = seconds - float64(minutes*60)
	duration := ""
	if hours > 0 {
		duration = duration + fmt.Sprintf("%d h ", hours)
	} else if minutes > 0 {
		duration = duration + fmt.Sprintf("%d m ", minutes)
	}
	return duration + fmt.Sprintf("%.2f s", seconds)
}

func SecondsToMilliseconds(seconds int64) int64 {
	return seconds * 1000
}

func MillisecondsToSeconds(milliseconds int64) int64 {
	return milliseconds / 1000
}

func ConvertUnixMilliToYearAndMonthFormat(milliseconds int64) string {
	_time := time.UnixMilli(milliseconds).UTC()
	date := time.Date(_time.Year(), _time.Month(), 1, 0, 0, 0, 0, time.UTC)
	return date.Format("2006-01")
}

func ConvertUnixToStringFormat(seconds int64) string {
	timestamp := time.Unix(seconds, 0).UTC()
	return timestamp.Format("2006-01-02 15:04:05")
}

func ConvertMilliUnixToStringFormat(milliseconds int64) string {
	timestamp := time.UnixMilli(milliseconds).UTC()
	return timestamp.Format("2006-01-02 15:04:05")
}

func ConvertMilliUnixToStringWithMillFormat(milliseconds int64) string {
	timestamp := time.UnixMilli(milliseconds).UTC()
	return timestamp.Format("2006-01-02 15:04:05.000")
}

func ParseDuration(d string) (time.Duration, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return 0, err
	}
	return duration, nil
}

package utils

import (
	"log"
	"strconv"
	"time"
)

type Time struct {
}

// UnixToTime 毫秒时间戳转时间
func (t *Time) UnixToTime(lastTime string) time.Time {
	timestamp, err := strconv.ParseInt(lastTime, 10, 64)
	log.Println("timestamp: ", timestamp)
	if err != nil {
		log.Println("strconv fail: ", err)
		return time.Time{}
	}
	tm := time.Unix(timestamp/1000, 0)
	return tm
}

// TimeToUnix 时间转毫秒时间戳
func (t *Time) TimeToUnix(e time.Time) int64 {
	timeUnix, _ := time.Parse("2006-01-02 15:04:05", e.Format("2006-01-02 15:04:05"))
	return timeUnix.UnixNano() / 1e6
}

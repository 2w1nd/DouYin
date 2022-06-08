package utils

import (
	"log"
	"strconv"
	"time"
)

// UnixToTime 毫秒时间戳转时间
func UnixToTime(lastTime string) time.Time {
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
func TimeToUnix(e time.Time) int64 {
	return e.UnixNano() / 1e6
}

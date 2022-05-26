package utils

import (
	"log"
	"strconv"
	"time"
)

type Time struct {
}

func (t *Time) GetTimeStamp(lastTime string) time.Time {
	timestamp, err := strconv.ParseInt(lastTime, 10, 64)
	log.Println("timestamp: ", timestamp)
	if err != nil {
		log.Println("strconv fail: ", err)
		return time.Time{}
	}
	tm := time.Unix(timestamp/1000, 0)
	return tm
}

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
	if err != nil {
		log.Println("strconv fail: ", err)
		return time.Time{}
	}
	tm := time.Unix(timestamp, 0)
	return tm
}

package initialize

import (
	"github.com/robfig/cron/v3"
)

func Cron() {
	c := cron.New()
	_, _ = c.AddFunc("@every 5s", func() {
		//fmt.Println("hello world")
	})
	c.Start()
}

package initialize

import (
	"github.com/DouYin/service/service"
	"github.com/robfig/cron/v3"
)

func Cron() {
	c := cron.New()
	_, _ = c.AddFunc("@every 5s", func() {
		service.SynchronizeDBAndRedis()
	})
	c.Start()
}

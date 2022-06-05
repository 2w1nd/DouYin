package initialize

import (
	"github.com/DouYin/service/service"
	"github.com/robfig/cron/v3"
)

/**
 * * /5 * * * * * 每5秒一次
 */

func Cron() {
	c := cron.New()
	// 每5分钟执行一次
	_, _ = c.AddFunc("*/1 * * * *", func() {
		service.SynchronizeDBAndRedis()
	})
	c.Start()
}

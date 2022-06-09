package initialize

import (
	"github.com/DouYin/service/cache"
	"github.com/robfig/cron/v3"
)

/**
 * * /5 * * * * * 每5秒一次
 */

func Cron() {
	c := cron.New()
	// 每5分钟同步一次关注粉丝关系
	_, _ = c.AddFunc("*/5 * * * *", func() {
		cache.SynchronizeRelationToDBFromRedis()
	})
	// 每10分钟同步一次视频点赞关系
	_, _ = c.AddFunc("*/10 * * * *", func() {
		cache.SynchronizeFavoriteDBFromRedis()
	})
	c.Start()
}

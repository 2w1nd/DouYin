package initialize

import (
	"github.com/DouYin/service/cache"
	"github.com/robfig/cron/v3"
)

/**
 * * /5 * * * * * 每5秒一次
 * @every 1s
 */

func Cron() {
	c := cron.New()
	// 每1个半小时同步一次关注粉丝关系
	_, _ = c.AddFunc("@every 1h30m", func() {
		cache.SynchronizeRelationToDBFromRedis()
	})
	// 每2小时同步一次视频点赞关系
	_, _ = c.AddFunc("@every 2h", func() {
		cache.SynchronizeFavoriteDBFromRedis()
	})
	c.Start()
}

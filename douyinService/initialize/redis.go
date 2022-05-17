package initialize

import (
	"context"
	"github.com/DouYin/service/global"
	"github.com/go-redis/redis/v8"
	"github.com/google/martian/log"
)

func Redis() {
	redisCfg := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Errorf("redis connect ping failed, err")
	} else {
		log.Infof("redis connect")
		global.REDIS = client
	}
}

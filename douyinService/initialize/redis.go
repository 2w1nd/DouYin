package initialize

import (
	"context"
	"github.com/DouYin/service/global"
	"github.com/go-redis/redis/v8"
	"github.com/google/martian/log"
	"time"
)

func Redis() *redis.Client {
	redisCfg := global.CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Errorf("redis connect ping failed", err)
	} else {
		log.Infof("redis connect")
		global.REDIS = client
	}
	return client
}

package cache

import (
	"context"
	"github.com/DouYin/service/global"
)

type UserCache struct {
}

func (uc *UserCache) GetUserNameById(id string) string {
	authorName, err := global.REDIS.HGet(context.Background(), "users:user", id).Result()
	if err != nil {
		return ""
	}
	return authorName
}

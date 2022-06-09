package cache

import (
	"context"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"strconv"
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

// GetAuthorNameInRedis 查作者名称
func (uc *UserCache) GetAuthorNameInRedis(authorId uint64) string {
	authorName, _ := global.REDIS.HGet(context.Background(), "users:user", strconv.FormatUint(authorId, 10)).Result()
	if authorName == "" { // 找不到，从数据库找
		where := model.User{UserId: authorId}
		author, _ := userRepository.GetFirstUser(where)
		authorName = author.Username
	}
	return authorName
}

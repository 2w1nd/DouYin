package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type UserRepository struct {
}

// QueryUserDtoInfo 查询UserDto需要的信息
func (u UserRepository) QueryUserDtoInfo(userId uint64) model.User {
	var user model.User
	query := global.DB.Debug().
		Model(model.User{})
	query.Select("id", "user_id", "name", "follow_count", "follower_count").
		Where("user_id = ?", userId).Limit(1)
	query.Find(&user)
	return user
}

// IsFollow 判断关注关系
// src 源 dst 目标
func (u UserRepository) IsFollow(src, dst uint64) bool {
	var count int64
	query := global.DB.Debug().
		Model(model.Follow{})
	query.
		Where("user_id = ? and followed_user_id = ? and is_deleted = 0", src, dst).Limit(1)
	query.Count(&count)
	return count > 0
}

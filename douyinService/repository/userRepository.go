package repository

import (
	"github.com/DouYin/common/constant"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type UserRepository struct {
	Base BaseRepository
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

func (u UserRepository) UpdateFollowCount(userId uint64, Type int) bool {
	var user model.User
	query := global.DB.Debug().
		Model(model.User{})
	query.Select("id", "user_id", "name", "follow_count", "follower_count").
		Where("user_id = ?", userId).Limit(1)
	query.Find(&user)
	db := global.DB.Where(user)
	var out model.User
	//var Cnt uint32
	//if cnt > 0 {
	//	Cnt = int32(cnt)
	//} else {
	//	Cnt = int32(-cnt)
	//}
	if Type == constant.FOCUS {
		if err := db.Model(out).Debug().Where(user).Update("follow_count", user.FollowCount+1).Error; err != nil {
			return false
		}
	} else if Type == constant.NoFOCUS {
		if err := db.Model(out).Debug().Where(user).Update("follow_count", user.FollowCount-1).Error; err != nil {
			return false
		}
	}
	return true
}

func (u UserRepository) UpdateFollowerCount(userId uint64, Type int) bool {
	var user model.User
	query := global.DB.Debug().
		Model(model.User{})
	query.Select("id", "user_id", "name", "follow_count", "follower_count").
		Where("user_id = ?", userId).Limit(1)
	query.Find(&user)
	db := global.DB.Where(user)
	var out model.User
	//var Cnt uint32
	//if cnt > 0 {
	//	Cnt = uint32(cnt)
	//} else {
	//	Cnt = uint32(-cnt)
	//}
	if Type == constant.FOCUS {
		if err := db.Model(out).Debug().Where(user).Update("follower_count", user.FollowerCount+1).Error; err != nil {
			return false
		}
	} else if Type == constant.NoFOCUS {
		if err := db.Model(out).Debug().Where(user).Update("follower_count", user.FollowerCount-1).Error; err != nil {
			return false
		}
	}
	return true
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

func (u UserRepository) CreateUser(user *model.User) (*model.User, error) {
	err := global.DB.Raw("INSERT INTO douyin_user(user_id,name,username,password,salt) VALUES(?,?,?,?,?);",
		user.UserId, user.Name, user.Username, user.Password, user.Salt).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetFirstUser(where interface{}) (model.User, error) {
	var user model.User
	if err := u.Base.First(where, &user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u UserRepository) GetUserByUserName(username string) (model.User, error) {
	user := model.User{}
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u UserRepository) GetfollowByUserId(userId, followerId uint64) (model.Follow, error) {
	follow := model.Follow{}
	err := global.DB.Where("user_id = ? AND followed_user_id = ?", userId, followerId).First(&follow).Error
	if err != nil {
		return model.Follow{}, err
	}
	return follow, nil
}

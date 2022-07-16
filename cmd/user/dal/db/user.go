package db

import (
	"context"
	"errors"
	"github.com/DouYin/pkg/constants"
	"gorm.io/gorm"
)

type User struct {
	Id            uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`     // 主键
	UserId        uint64 `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                      // 用户id，唯一
	Name          string `gorm:"column:name;type:varchar(32)" json:"name"`                                    // 用户昵称，业务层长度限制
	Username      string `gorm:"column:username;type:varchar(255)" json:"username"`                           // 登录名，长度限制没确定
	Password      string `gorm:"column:password;type:varchar(255)" json:"password"`                           // 加密后密码
	Salt          string `gorm:"column:salt;type:varchar(22)" json:"salt"`                                    // 随机盐
	FollowCount   uint32 `gorm:"column:follow_count;type:int(10) unsigned;default:0" json:"follow_count"`     // 关注数，冗余字段
	FollowerCount uint32 `gorm:"column:follower_count;type:int(10) unsigned;default:0" json:"follower_count"` // 粉丝数，冗余字段
	gorm.Model
}

func (u *User) TableName() string {
	return constants.UserTableName
}

// Create 新建用户，将新记录插入数据库
func Create(ctx context.Context, user *User) (*User, error) {
	if err := DB.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// MGet 批量获取用户信息
func MGet(ctx context.Context, ids []int64) ([]*User, error) {
	res := make([]*User, 0)
	if len(ids) == 0 {
		return res, nil
	}

	if err := DB.WithContext(ctx).Where("id in ?", ids).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

// QueryUser 根据用户名获取用户信息
func QueryUser(ctx context.Context, username string) (*User, error) {
	user := &User{}
	if err := DB.WithContext(ctx).Where("user_name = ?", username).Take(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func IsExist(ctx context.Context, id int64) (bool, error) {
	res := &User{}
	err := DB.WithContext(ctx).Where("id = ?", id).Take(&res).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return true, nil
}

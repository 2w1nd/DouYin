package model

type User struct {
	Id            uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`     // 主键
	UserId        uint64 `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                      // 用户id，唯一
	Name          string `gorm:"column:name;type:varchar(32)" json:"name"`                                    // 用户昵称，业务层长度限制
	Username      string `gorm:"column:username;type:varchar(255)" json:"username"`                           // 登录名，长度限制没确定
	Password      string `gorm:"column:password;type:varchar(255)" json:"password"`                           // 加密后密码
	Salt          string `gorm:"column:salt;type:varchar(22)" json:"salt"`                                    // 随机盐
	FollowCount   uint32 `gorm:"column:follow_count;type:int(10) unsigned;default:0" json:"follow_count"`     // 关注数，冗余字段
	FollowerCount uint32 `gorm:"column:follower_count;type:int(10) unsigned;default:0" json:"follower_count"` // 粉丝数，冗余字段
	BaseModel
}

func (m *User) TableName() string {
	return "douyin_user"
}

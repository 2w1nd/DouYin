package model

type Follow struct {
	Id             uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`  // 主键
	UserId         uint64 `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                   // 关注者id
	FollowedUserId uint64 `gorm:"column:followed_user_id;type:bigint(20) unsigned" json:"followed_user_id"` // 被关注者id
	IsDeleted      bool   `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted"`                      // 关注状态，0关注，1取消
	BaseModel
}

func (m *Follow) TableName() string {
	return "douyin_follow"
}

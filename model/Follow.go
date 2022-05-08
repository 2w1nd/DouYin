package model

import "database/sql"

type Follow struct {
	Id             uint64         `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`         // 主键
	UserId         sql.NullInt64  `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                          // 关注者id
	FollowedUserId sql.NullInt64  `gorm:"column:followed_user_id;type:bigint(20) unsigned" json:"followed_user_id"`        // 被关注者id
	IsDeleted      sql.NullInt32  `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted"`                             // 关注状态，0关注，1取消
	GmtModified    sql.NullTime   `gorm:"column:gmt_modified;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_modified"` // 更新时间
	GmtCreated     sql.NullTime   `gorm:"column:gmt_created;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_created"`   // 创建时间
	Ext            sql.NullString `gorm:"column:ext;type:varchar(255)" json:"ext"`                                         // 扩展字段
}

func (m *Follow) TableName() string {
	return "douyin_follow"
}

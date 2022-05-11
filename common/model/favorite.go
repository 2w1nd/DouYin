package model

import "database/sql"

type Favorite struct {
	Id        uint64        `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"` // 主键
	UserId    sql.NullInt64 `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                  // 用户id
	VideoId   sql.NullInt64 `gorm:"column:video_id;type:bigint(20) unsigned" json:"video_id"`                // 视频id
	IsDeleted sql.NullInt32 `gorm:"column:is_deleted;type:tinyint(1)" json:"is_deleted"`                     // 点赞状态，0点赞，1取消
	BaseModel
}

func (m *Favorite) TableName() string {
	return "douyin_favorite"
}

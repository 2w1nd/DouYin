package model

import (
	"time"
)

type Video struct {
	Id            uint64    `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`         // 主键
	VideoId       uint64    `gorm:"column:video_id;type:bigint(20) unsigned" json:"video_id"`                        // 唯一videoid
	AuthorId      uint64    `gorm:"column:author_id;type:bigint(20) unsigned" json:"author_id"`                      // 作者id
	Description   string    `gorm:"column:description;type:varchar(255)" json:"description"`                         // 视频描述
	Path          string    `gorm:"column:path;type:varchar(255)" json:"path"`                                       // 视频存储路径
	CoverPath     string    `gorm:"column:cover_path;type:varchar(255)" json:"cover_path"`                           // 视频封面路径
	FavoriteCount uint32    `gorm:"column:favorite_count;type:int(10) unsigned;default:0" json:"favorite_count"`     // 点赞数，冗余字段
	CommentCount  uint32    `gorm:"column:comment_count;type:int(10) unsigned;default:0" json:"comment_count"`       // 评论数，冗余字段
	GmtModified   time.Time `gorm:"column:gmt_modified;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_modified"` // 更新时间
	GmtCreated    time.Time `gorm:"column:gmt_created;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_created"`   // 创建时间
	Ext           string    `gorm:"column:ext;type:varchar(255)" json:"ext"`                                         // 扩展字段
}

func (m *Video) TableName() string {
	return "douyin_video"
}

package model

import "database/sql"

type Comment struct {
	Id          uint64         `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"`         // 主键
	CommentId   sql.NullInt64  `gorm:"column:comment_id;type:bigint(20) unsigned" json:"comment_id"`                    // 唯一commentId
	UserId      sql.NullInt64  `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                          // 评论用户id
	VideoId     sql.NullInt64  `gorm:"column:video_id;type:bigint(20) unsigned" json:"video_id"`                        // 视频id
	Content     sql.NullString `gorm:"column:content;type:varchar(255)" json:"content"`                                 // 评论内容
	Pid         sql.NullInt64  `gorm:"column:pid;type:bigint(20) unsigned" json:"pid"`                                  // 父id，多级评论用
	GmtModified sql.NullTime   `gorm:"column:gmt_modified;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_modified"` // 更新时间
	GmtCreated  sql.NullTime   `gorm:"column:gmt_created;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_created"`   // 创建时间
	Ext         sql.NullString `gorm:"column:ext;type:varchar(255)" json:"ext"`                                         // 扩展字段
}

func (m *Comment) TableName() string {
	return "douyin_comment"
}

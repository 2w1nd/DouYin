package model

import "database/sql"

type Comment struct {
	Id        uint64         `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"` // 主键
	CommentId sql.NullInt64  `gorm:"column:comment_id;type:bigint(20) unsigned" json:"comment_id"`            // 唯一commentId
	UserId    sql.NullInt64  `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                  // 评论用户id
	VideoId   sql.NullInt64  `gorm:"column:video_id;type:bigint(20) unsigned" json:"video_id"`                // 视频id
	Content   sql.NullString `gorm:"column:content;type:varchar(255)" json:"content"`                         // 评论内容
	Pid       sql.NullInt64  `gorm:"column:pid;type:bigint(20) unsigned" json:"pid"`                          // 父id，多级评论用
	BaseModel
}

func (m *Comment) TableName() string {
	return "douyin_comment"
}

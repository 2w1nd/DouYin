package model

type Comment struct {
	Id        uint64         `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" json:"id"` // 主键
	CommentId uint64 `gorm:"column:comment_id;type:bigint(20) unsigned" json:"comment_id"`            // 唯一commentId
	UserId    uint64  `gorm:"column:user_id;type:bigint(20) unsigned" json:"user_id"`                  // 评论用户id
	VideoId   uint64  `gorm:"column:video_id;type:bigint(20) unsigned" json:"video_id"`                // 视频id
	Content   string `gorm:"column:content;type:varchar(255)" json:"content"`                         // 评论内容
	Pid       uint64  `gorm:"column:pid;type:bigint(20) unsigned" json:"pid"`                          // 父id，多级评论用
	CommentUser []User `json:"comment_user" gorm:"foreignKey:UserId;references:UserId"` // 用于预加载，每个user可以有许多Comment
	BaseModel
}

func (m *Comment) TableName() string {
	return "douyin_comment"
}

package model

type Comment struct {
	Id         int64  `json:"id,omitempty" form:"id" gorm:"comment:评论ID"`
	UserID     int64  `json:"user_id"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

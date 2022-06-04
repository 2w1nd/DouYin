package model

type Favorite struct {
	Id        uint64 `gorm:"column:id;type:bigint(20) unsigned;primary_key;AUTO_INCREMENT" ` // 主键
	UserId    uint64 `gorm:"column:user_id;type:bigint(20) unsigned" `                       // 用户id
	VideoId   uint64 `gorm:"column:video_id;type:bigint(20) unsigned" `                      // 视频id
	IsDeleted bool   `gorm:"column:is_deleted;type:tinyint(1);default:true" `                // 点赞状态，0点赞，1取消
	User      User   `json:"user" gorm:"foreignKey:UserId;references:UserId"`                // 用于预加载
	Video     Video  `json:"video" gorm:"foreignKey:VideoId;references:VideoId"`             // 用于预加载
	BaseModel
}

func (m *Favorite) TableName() string {
	return "douyin_favorite"
}

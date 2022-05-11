package model

import (
	"time"
)

type BaseModel struct {
	GmtModified time.Time   `gorm:"column:gmt_modified;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_modified"` // 更新时间
	GmtCreated  time.Time   `gorm:"column:gmt_created;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_created"`   // 创建时间
	Ext         string `gorm:"column:ext;type:varchar(255)" json:"ext"`                                         // 扩展字段
}

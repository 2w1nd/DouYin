package model

import "database/sql"

type BaseModel struct {
	GmtModified sql.NullTime   `gorm:"column:gmt_modified;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_modified"` // 更新时间
	GmtCreated  sql.NullTime   `gorm:"column:gmt_created;type:datetime;default:CURRENT_TIMESTAMP" json:"gmt_created"`   // 创建时间
	Ext         sql.NullString `gorm:"column:ext;type:varchar(255)" json:"ext"`                                         // 扩展字段
}

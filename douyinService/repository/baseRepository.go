package repository

import "github.com/DouYin/service/global"

type BaseRepository struct {
}

// Create
// @Description: 创建实体
// @receiver: b
// @param: value
// @return: error
func (b *BaseRepository) Create(value interface{}) error {
	return global.DB.Create(value).Error
}

// Save
// @Description: 删除实体
// @receiver: b
// @param: value
// @return: error
func (b *BaseRepository) Save(value interface{}) error {
	return global.DB.Save(value).Error
}

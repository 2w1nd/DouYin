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

// DeleteByID
// @Description: 根据id删除实体
// @receiver: b
// @param: model
// @param: id
// @return: error
func (b *BaseRepository) DeleteByID(model interface{}, id int) error {
	return global.DB.Where("id = ?", id).Delete(model).Error
}

// First
// @Description: 根据条件获取一个实体
// @receiver: b
// @param: where
// @param: out
// @param: selects
// @return: error
func (b *BaseRepository) First(where interface{}, out interface{}, selects ...string) error {
	db := global.DB.Where(where)
	if len(selects) > 0 {
		for _, sel := range selects {
			db = db.Select(sel)
		}
	}
	return db.First(out).Error
}

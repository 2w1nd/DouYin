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

// GetPages
// @Description: 分页返回数据
// @receiver: b
// @param: model
// @param: out
// @param: pageIndex
// @param: pageSize
// @param: totalCount
// @param: where
// @param: orders
// @return: error
func (b *BaseRepository) GetPages(model interface{}, out interface{}, pageIndex, pageSize int, totalCount *int64, where interface{}, orders ...string) error {
	db := global.DB.Model(model).Where(model)
	db = db.Where(where)
	if len(orders) > 0 {
		for _, order := range orders {
			db = db.Order(order)
		}
	}
	err := db.Count(totalCount).Error
	if err != nil {
		return err
	}
	if *totalCount == 0 {
		return nil
	}
	return db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Error
}

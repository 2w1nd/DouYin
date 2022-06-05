package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"gorm.io/gorm/clause"
	"log"
)

type FavoriteRepository struct {
	Base BaseRepository
}

func (fr *FavoriteRepository) GetFavoriteByUserIdAndVideoId(userId uint64, videoId uint64) (bool, model.Favorite) {
	var count int64
	var favorite model.Favorite
	query := global.DB.Model(model.Favorite{}).Where("user_id=? and video_id=?", userId, videoId)
	query.Count(&count)
	if count >= 1 {
		query.Find(&favorite)
		return true, favorite
	} else {
		return false, model.Favorite{}
	}
}
func (fr *FavoriteRepository) AddFavorite(favorite model.Favorite) bool {
	log.Println("添加到数据库")
	if err := global.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_deleted"}),
	}).Create(&favorite); err != nil {
		return false
	}
	return true
}
func (fr *FavoriteRepository) DeleteFavoriteById(where interface{}) bool {
	var favorite model.Favorite
	if err := fr.Base.DeleteByID(favorite, where); err != nil {
		return false
	}
	return true
}

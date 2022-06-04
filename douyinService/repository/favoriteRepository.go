package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type FavoriteRepository struct {
	Base BaseRepository
}

func (fr *FavoriteRepository) GetFavoriteByUserIdAndVideoId(userId int64, videoId int64) (bool, model.Favorite) {
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
	if err := fr.Base.Create(&favorite); err != nil {
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

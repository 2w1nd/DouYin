package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"gorm.io/gorm/clause"
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
	//log.Println("添加到数据库: is_del: ", favorite.IsDeleted)
	if err := global.DB.Debug().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"is_deleted"}),
	}).Create(&favorite); err != nil {
		return false
	}
	return true

	//var favoriteInfo model.Favorite
	//result := global.DB.Model(model.Favorite{}).Where("user_id =? and video_id =?", favorite.UserId, favorite.VideoId).First(&favoriteInfo)
	//if result.Error != nil {
	//	log.Println(result.Error)
	//	if result.Error == gorm.ErrRecordNotFound {
	//		result := global.DB.Model(model.Favorite{}).Create(&favorite)
	//		if result.Error != nil {
	//			return false
	//		}
	//	}
	//} else {
	//	result := global.DB.Debug().Model(model.Favorite{}).Where("user_id=? and video_id=?", favorite.UserId, favorite.VideoId).Update("is_deleted", favorite.IsDeleted)
	//	if result.Error != nil {
	//		return false
	//	}
	//}
	//return true
}
func (fr *FavoriteRepository) DeleteFavoriteById(where interface{}) bool {
	var favorite model.Favorite
	if err := fr.Base.DeleteByID(favorite, where); err != nil {
		return false
	}
	return true
}

// SaveFavoriteCount 根据VideoId查询对应Video点赞数量，并存入数据库
func (fr *FavoriteRepository) SaveFavoriteCount(videoId uint64, ans int64) error {
	//id, err := strconv.ParseInt(videoId, 10, 64)
	//ans, err := favoriteService.RedisGetVideoFavoriteCount(id)
	global.DB.Model(model.Video{}).Where("video_id = ?", videoId).Update("favorite_count", ans)
	return nil
}

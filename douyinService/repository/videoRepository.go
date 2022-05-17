package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type VideoRepository struct {
}

// GetVideoWithAuthor
// @Description: 获得视频以及对应的作者信息
// @receiver: v
// @return: []model.Video
func (v *VideoRepository) GetVideoWithAuthor(latestTime string) []model.Video {
	var videoList []model.Video
	query := global.DB.Model(model.Video{}).Preload("User")
	//if latestTime != "" {
	//	query.Where("gmt_created <= ?", latestTime)
	//}
	query.Order("gmt_created DESC").Limit(30)
	query.Debug().Find(&videoList)
	return videoList
}

// GetVideoWithAuthorAndFollowAndFavorite
// @Description: 查出视频以及对应作者信息，以及当前登录用户是否follow了该作者，是否favorite了该视频
// @receiver: v
func (v *VideoRepository) GetVideoWithAuthorAndFollowAndFavorite(latestTime string, id uint64) []model.Video {
	var videoList []model.Video
	query := global.DB.Debug().
		Model(model.Video{}).
		Preload("User").
		Preload("User.FollowedUser").
		Preload("Favorite", "user_id = ?", id)

	//if latestTime != "" {
	//	query.Where("gmt_created <= ?", latestTime)
	//}
	query.Order("gmt_created DESC").Limit(30)
	query.Find(&videoList)
	return videoList
}

package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"log"
	"time"
)

type VideoRepository struct {
}

// GetVideoWithAuthor
// @Description: 获得视频以及对应的作者信息
// @receiver: v
// @return: []model.Video
func (v *VideoRepository) GetVideoWithAuthor(latestTime time.Time) []model.Video {
	var videoList []model.Video
	query := global.DB.Model(model.Video{}).Preload("User")
	log.Println(latestTime)
	if !latestTime.IsZero() {
		query.Where("gmt_created <= ?", latestTime)
	}
	query.Order("gmt_created DESC").Limit(30)
	query.Debug().Find(&videoList)
	return videoList
}

// GetVideoWithAuthorAndFollowAndFavorite
// @Description: 查出视频以及对应作者信息，以及当前登录用户是否follow了该作者，是否favorite了该视频
// @receiver: v
func (v *VideoRepository) GetVideoWithAuthorAndFollowAndFavorite(latestTime time.Time, id uint64) []model.Video {
	var videoList []model.Video
	query := global.DB.Debug().
		Model(model.Video{}).
		Preload("User").
		Preload("User.FollowedUser").
		Preload("Favorite", "user_id = ?", id)

	if !latestTime.IsZero() {
		query.Where("gmt_created <= ?", latestTime)
	}
	query.Order("gmt_created DESC").Limit(30)
	query.Find(&videoList)
	return videoList
}

// GetPublishList
// 获取用户发布的视频列表,不填充author和Favorite，pageIndex从1开始
func (v *VideoRepository) GetPublishList(userId uint64, pageIndex, pageSize int) []model.Video {
	//查询1652259488发布的视频，并且判断995有没有点过赞
	//SELECT t1.*,(t2.id is NULL) as is_favorite from douyin_video as t1 left outer join douyin_favorite as t2
	//ON t1.video_id = t2.video_id and t2.user_id = '995' where t1.author_id = '1652259488'
	var videoList []model.Video
	query := global.DB.Debug().
		Model(model.Video{})
	query.Where("author_id = ?", userId).Offset(pageIndex - 1*pageSize).Limit(pageSize).Order("gmt_created DESC")
	query.Find(&videoList)
	return videoList
}

// GetPublishListWithFavorite
// 获取用户发布的视频列表,填充和Favorite，不填充author，pageIndex从1开始
func (v *VideoRepository) GetPublishListWithFavorite(userId uint64, pageIndex, pageSize int, loginUser uint64) []model.Video {
	//查询1652259488发布的视频，并且判断995有没有点过赞
	//SELECT t1.*,(t2.id is NULL) as is_favorite from douyin_video as t1 left outer join douyin_favorite as t2
	//ON t1.video_id = t2.video_id and t2.user_id = '995' where t1.author_id = '1652259488'
	var videoList []model.Video
	query := global.DB.Debug().
		Model(model.Video{}).
		Preload("Favorite", "user_id = ?", loginUser)
	query.Where("author_id = ?", userId).Offset(pageIndex - 1*pageSize).Limit(pageSize).Order("gmt_created DESC")
	query.Find(&videoList)
	return videoList
}

func (v *VideoRepository) SaveVideo(video model.Video) uint64 {
	result := global.DB.Debug().Create(&video)
	err := result.Error
	if err != nil {
		return 0
	}
	return video.Id

}

func (v *VideoRepository) GetVideoByVideoId(videoId uint64) model.Video {
	video := model.Video{}
	query := global.DB.Debug().Model(model.Video{}).Preload("User").Where("video_id = ?", videoId)
	query.Find(&video)
	return video
}

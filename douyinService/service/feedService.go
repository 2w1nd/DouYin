package service

import (
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/cache"
	"github.com/DouYin/service/utils"
	"time"
)

type FeedService struct {
	videoCache cache.VideoCache
}

// Feed 视频流接口
func (fs *FeedService) Feed(userId uint64, latestTime string) ([]vo.VideoVo, time.Time) {
	var (
		videoVos  []vo.VideoVo
		videoList []model.Video
		nextTime  int64
	)
	//从缓存查询
	videoVos, nextTime = fs.videoCache.ReadFeedDataFromRedis(userId)
	if len(videoVos) != 0 {
		return videoVos, time.Unix(nextTime/1000, 0)
	}

	// 从数据库中查询
	if userId == 0 {
		videoList = videoRepository.GetVideoWithAuthor(utils.UnixToTime(latestTime))
	} else {
		videoList = videoRepository.GetVideoWithAuthorAndFollowAndFavorite(utils.UnixToTime(latestTime), userId)
	}
	if len(videoList) == 0 {
		return []vo.VideoVo{}, time.Time{}
	}

	videoVos = fs.videoList2Vo(videoList)
	// 放入缓存
	fs.videoCache.LoadFeedDataToRedis(videoList)
	return videoVos, videoList[0].GmtCreated
}

// videoList2Vo 将查出来的数据传入vo
func (fs *FeedService) videoList2Vo(videoList []model.Video) []vo.VideoVo {
	var videoVos []vo.VideoVo
	for _, video := range videoList {
		var isFollow, isFavorite bool
		if len(video.User.FollowedUser) != 0 {
			isFollow = video.User.FollowedUser[0].IsDeleted
		} else {
			isFollow = false
		}
		if len(video.Favorite) != 0 {
			isFavorite = false
		} else {
			isFavorite = false
		}
		videoVo := vo.VideoVo{
			VideoID: video.VideoId,
			Author: vo.AuthorVo{
				UserID:        video.User.UserId,
				Name:          video.User.Username,
				FollowCount:   video.User.FollowCount,
				FollowerCount: video.User.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       video.Path,
			CoverUrl:      video.CoverPath,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			Title:         video.Title,
			IsFavorite:    isFavorite,
		}
		videoVos = append(videoVos, videoVo)
	}
	return videoVos
}

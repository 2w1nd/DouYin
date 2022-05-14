package service

import (
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/repository"
	"log"
)

type FeedService struct {
}

var videoRepository repository.VideoRepository

// Feed
// @Description: 视频流接口
// @receiver: fs
// @param: token
// @param: latestTime
// @return: []vo.VideoVo
func (fs *FeedService) Feed(token, latestTime string) []vo.VideoVo {
	var (
		videoVos  []vo.VideoVo
		videoList []model.Video
	)
	// TODO 解析token获取ID
	if token == "" {
		videoList = videoRepository.GetVideoWithAuthor(latestTime)
	} else {
		videoList = videoRepository.GetVideoWithAuthorAndFollowAndFavorite(latestTime, 90071992547409929)
	}
	videoVos = fs.videoList2Vo(videoList)
	log.Println(videoVos)
	return videoVos
}

//
// @Description: 将查出来的数据传入vo
// @receiver: fs
// @param: videoList
// @return: []vo.VideoVo
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
			isFavorite = video.Favorite[0].IsDeleted
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
			IsFavorite:    isFavorite,
		}
		videoVos = append(videoVos, videoVo)
	}
	return videoVos
}

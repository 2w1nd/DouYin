package service

import (
	"encoding/json"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"golang.org/x/net/context"
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
func (fs *FeedService) Feed(id uint64, latestTime string) []vo.VideoVo {
	var (
		videoVos  []vo.VideoVo
		videoList []model.Video
	)

	// 从缓存查询
	data1, _ := global.REDIS.Get(context.Background(), "videoVos").Result()
	if data1 != "" {
		err := json.Unmarshal([]byte(data1), &videoVos)
		if err != nil {
			return nil
		}
		log.Println(videoVos)
		if len(videoVos) != 0 {
			return videoVos
		}
	}
	log.Println("从数据库中查询")
	if id == 0 {
		videoList = videoRepository.GetVideoWithAuthor(latestTime)
	} else {
		videoList = videoRepository.GetVideoWithAuthorAndFollowAndFavorite(latestTime, 90071992547409929)
	}
	if len(videoList) == 0 {
		return []vo.VideoVo{}
	}

	videoVos = fs.videoList2Vo(videoList)
	// 放入缓存
	data, _ := json.Marshal(videoVos)
	global.REDIS.Set(context.Background(), "videoVos", data, 0)
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

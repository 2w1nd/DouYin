package service

import (
	"encoding/json"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"golang.org/x/net/context"
	"log"
	"time"
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
func (fs *FeedService) Feed(id uint64, latestTime string) ([]vo.VideoVo, time.Time) {
	var (
		videoData vo.VideoData
		videoVos  []vo.VideoVo
		videoList []model.Video
		timeUtil  utils.Time
	)
	// 从缓存查询
	data1, _ := global.REDIS.Get(context.Background(), "videoVos").Result()
	if data1 != "" {
		log.Println("从缓存中查询")
		err := json.Unmarshal([]byte(data1), &videoData)
		if err != nil {
			return nil, time.Time{}
		}
		if len(videoData.VideoList) != 0 {
			return videoData.VideoList, time.Unix(videoData.NextTime, 0)
		}
	}
	// 从数据库中查询
	log.Println("从数据库中查询")
	if id == 0 {
		videoList = videoRepository.GetVideoWithAuthor(timeUtil.UnixToTime(latestTime))
	} else {
		videoList = videoRepository.GetVideoWithAuthorAndFollowAndFavorite(timeUtil.UnixToTime(latestTime), id)
	}
	if len(videoList) == 0 {
		return []vo.VideoVo{}, time.Time{}
	}

	videoVos = fs.videoList2Vo(videoList)
	// 放入缓存
	videoData.VideoList = videoVos
	videoData.NextTime = timeUtil.TimeToUnix(videoList[0].GmtCreated)
	data, _ := json.Marshal(videoData)
	global.REDIS.Set(context.Background(), "videoVos", data, 10*time.Minute)
	return videoVos, videoList[0].GmtCreated
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

package cache

import (
	"context"
	"encoding/json"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"log"
	"strconv"
)

type VideoCache struct {
	favoriteCache FavoriteCache
	commentCache  CommentCache
	userCache     UserCache
	relationCache RelationCache
}

type VideoMsg struct {
	VideoID    uint64 `json:"id,omitempty"`
	AuthorID   uint64 `json:"author_id,omitempty"`
	PlayUrl    string `json:"play_url,omitempty"`
	CoverUrl   string `json:"cover_url,omitempty"`
	Title      string `json:"title,omitempty"`
	CreateTime int64  `json:"create_time"`
}

var userRepository repository.UserRepository
var videoRepository repository.VideoRepository
var relationCache RelationCache

// ReadFeedDataFromRedis 从Redis中读取视频数据，若读取不到，则去DB找
func (vc *VideoCache) ReadFeedDataFromRedis(userId uint64) (videoVos []vo.VideoVo, nextTime int64) {
	// 查视频数据
	videosIds, _ := global.REDIS.LRange(context.Background(), "videoIds", 0, -1).Result()
	videoVos, nextTime = vc.GetVideoVoByIdsFromRedis(videosIds, userId)
	return videoVos, nextTime
}

// LoadFeedDataToRedis 加载视频流数据到Redis中
func (vc *VideoCache) LoadFeedDataToRedis(videoModels []model.Video) {
	for _, video := range videoModels {
		videomsg := VideoMsg{
			VideoID:    video.VideoId,
			AuthorID:   video.AuthorId,
			PlayUrl:    video.Path,
			CoverUrl:   video.CoverPath,
			Title:      video.Title,
			CreateTime: utils.TimeToUnix(video.GmtCreated),
		}
		videoMsgJson, _ := json.Marshal(videomsg)
		str := global.REDIS.HGet(context.Background(), "videos", strconv.FormatUint(video.VideoId, 10)).Err()
		if str != nil {
			global.REDIS.HMSet(context.Background(), "videos", strconv.FormatUint(video.VideoId, 10), videoMsgJson)
			global.REDIS.LPush(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10))
		}
	}
}

// ReadPublishDataFromRedis 从Redis中读取数据，若读取不到，则去DB中找
func (vc *VideoCache) ReadPublishDataFromRedis(userId, myId uint64) []vo.VideoVo {
	// 得到该用户得视频id列表
	videoIdsStr, _ := global.REDIS.LRange(context.Background(), "userVideos:userVideo"+strconv.FormatUint(userId, 10), 0, -1).Result()
	videoVos, _ := vc.GetVideoVoByIdsFromRedis(videoIdsStr, myId)
	return videoVos
}

// LoadPublishDataToRedis 加载DB数据到Redis中
func (vc *VideoCache) LoadPublishDataToRedis(video model.Video) {
	videomsg := VideoMsg{
		VideoID:    video.VideoId,
		AuthorID:   video.AuthorId,
		PlayUrl:    video.Path,
		CoverUrl:   video.CoverPath,
		Title:      video.Title,
		CreateTime: utils.TimeToUnix(video.GmtCreated),
	}
	videoMsgJson, _ := json.Marshal(videomsg)
	str := global.REDIS.HGet(context.Background(), "videos", strconv.FormatUint(video.VideoId, 10)).Err()
	if str != nil {
		global.REDIS.HMSet(context.Background(), "videos", strconv.FormatUint(video.VideoId, 10), videoMsgJson)
		global.REDIS.LPush(context.Background(), "userVideos:userVideo"+strconv.FormatUint(video.AuthorId, 10), strconv.FormatUint(video.VideoId, 10))
		global.REDIS.LPush(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10))
	}
}

func (vc *VideoCache) GetVideoVoByIdsFromRedis(videosIds []string, userId uint64) (videoVos []vo.VideoVo, nextTime int64) {
	// 查视频数据
	var videoMsgs []VideoMsg
	for _, videoId := range videosIds {
		video := global.REDIS.HGet(context.Background(), "videos", videoId).String()
		videom := utils.SplitStringForList(video, ":")
		var v VideoMsg
		err := json.Unmarshal([]byte(videom), &v)
		if err != nil {
			log.Println(err)
		}
		videoMsgs = append(videoMsgs, v)
	}
	if len(videoMsgs) == 0 { // 如果找不到视频流数据，说明已经过期，去DB中找
		return
	}
	for _, videoMsg := range videoMsgs {
		// 查作者名称
		authorName := vc.userCache.GetAuthorNameInRedis(videoMsg.AuthorID)
		// 查粉丝数量，关注数量，当前用户是否关注
		followerCount, followCount, isFollow := vc.relationCache.GetFollowerCountAndFollowCountAndIsFollow(userId, videoMsg.AuthorID)
		// 查点赞数量，当前用户是否点赞
		favoriteCount, isFavorite := vc.favoriteCache.GetFavoriteCountAndIsFavorite(userId, videoMsg.VideoID)
		// 查评论数量
		commentCount := vc.commentCache.GetCommentCount(videoMsg.VideoID)

		videoVo := vo.VideoVo{
			VideoID: videoMsg.VideoID,
			Author: vo.AuthorVo{
				UserID:        videoMsg.AuthorID,
				Name:          authorName,
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       videoMsg.PlayUrl,
			CoverUrl:      videoMsg.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         videoMsg.Title,
		}
		videoVos = append(videoVos, videoVo)
	}
	if len(videoMsgs) >= 1 {
		nextTime = videoMsgs[len(videoMsgs)-1].CreateTime
	}
	return
}

package cache

import (
	"context"
	"encoding/json"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/entity/dto"
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
	var videoMsgs []VideoMsg
	videosIds, _ := global.REDIS.LRange(context.Background(), "videoIds", 0, -1).Result()
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
		authorName := vc.getAuthorNameInRedis(videoMsg.AuthorID)
		// 查粉丝数量，关注数量，当前用户是否关注
		followerCount, followCount, isFollow := vc.getFollowCountAndFollowedCountAndIsFollow(userId, videoMsg.AuthorID)
		// 查点赞数量，当前用户是否点赞
		favoriteCount, isFavorite := vc.getFavoriteCountAndIsFavorite(userId, videoMsg.VideoID)
		log.Println("feedService：favoriteCount isFavorite", videoMsg.Title, favoriteCount, isFavorite)
		// 查评论数量
		commentCount := vc.getCommentCount(videoMsg.VideoID)

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
			global.REDIS.RPush(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10))
		}
	}
}

// ReadPublishDataFromRedis 从Redis中读取数据，若读取不到，则去DB中找
func (vc *VideoCache) ReadPublishDataFromRedis(userId, myId uint64) []dto.VideoDto {
	// 得到该用户得视频id列表
	videoIdsStr, _ := global.REDIS.LRange(context.Background(), "userVideos:userVideo"+strconv.FormatUint(userId, 10), 0, -1).Result()
	// 遍历视频id
	var videoMsgs []VideoMsg
	var videoVos []dto.VideoDto
	for _, videoId := range videoIdsStr {
		video := global.REDIS.HGet(context.Background(), "videos", videoId).String()
		videoMsg := utils.SplitStringForList(video, ":")
		var v VideoMsg
		err := json.Unmarshal([]byte(videoMsg), &v)
		if err != nil {
			log.Println(err)
		}
		videoMsgs = append(videoMsgs, v)
	}
	if len(videoMsgs) == 0 { // 如果找不到视频流数据，说明已经过期，去DB中找
		return videoVos
	}
	for _, videoMsg := range videoMsgs {
		// 查作者名称
		authorName := vc.getAuthorNameInRedis(videoMsg.AuthorID)
		// 查粉丝数量，关注数量，当前用户是否关注
		followerCount, followCount, isFollow := vc.getFollowCountAndFollowedCountAndIsFollow(myId, userId)
		// 查点赞数量，当前用户是否点赞
		favoriteCount, isFavorite := vc.getFavoriteCountAndIsFavorite(myId, videoMsg.VideoID)
		// 查评论数量
		commentCount := vc.getCommentCount(videoMsg.VideoID)

		videoVo := dto.VideoDto{
			Id: videoMsg.VideoID,
			Author: dto.UserDto{
				Id:            videoMsg.AuthorID,
				Name:          authorName,
				FollowCount:   followCount,
				FollowerCount: followerCount,
				IsFollow:      isFollow,
			},
			PlayURL:       videoMsg.PlayUrl,
			CoverURL:      videoMsg.CoverUrl,
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Title:         videoMsg.Title,
		}
		videoVos = append(videoVos, videoVo)
	}
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
	str := global.REDIS.HGet(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10)).Err()
	if str != nil {
		global.REDIS.RPush(context.Background(), "video", videoMsgJson)
		global.REDIS.RPush(context.Background(), "userVideos:userVideo"+strconv.FormatUint(video.AuthorId, 10), strconv.FormatUint(video.VideoId, 10))
		global.REDIS.HMSet(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10), strconv.FormatUint(video.AuthorId, 10))
	}
}

// getAuthorNameInRedis 查作者名称
func (vc *VideoCache) getAuthorNameInRedis(authorId uint64) string {
	authorName, _ := global.REDIS.HGet(context.Background(), "users:user", strconv.FormatUint(authorId, 10)).Result()
	if authorName == "" { // 找不到，从数据库找
		where := model.User{UserId: authorId}
		author, _ := userRepository.GetFirstUser(where)
		authorName = author.Username
	}
	return authorName
}

// getFollowCountAndFollowedCountAndIsFollow 查粉丝数量，关注数量，当前用户是否关注
func (vc *VideoCache) getFollowCountAndFollowedCountAndIsFollow(myId uint64, authorId uint64) (uint32, uint32, bool) {
	var (
		followerCount, followCount uint32
		isFollow                   bool
		code1, code2, code3        int
	)
	followerCount, code1 = relationCache.RedisGetFollowerCount(int64(authorId))
	followCount, code2 = relationCache.RedisGetFollowCount(int64(authorId))
	if code1 == codes.RedisNotFound || code2 == codes.RedisNotFound { // 查DB
		where := model.User{UserId: authorId}
		author, _ := userRepository.GetFirstUser(where)
		followerCount, followCount = author.FollowerCount, author.FollowCount
	}
	isFollow, code3 = relationCache.RedisIsRelationCreated(int64(myId), int64(authorId))
	if code3 == codes.RedisNotFound {
		isFollow = userRepository.IsFollow(myId, authorId) // 查DB
	}
	return followerCount, followCount, isFollow
}

// getFavoriteCountAndIsFavorite 查点赞数量，当前用户是否点赞
func (vc *VideoCache) getFavoriteCountAndIsFavorite(userId uint64, videoId uint64) (uint32, bool) {
	isFavorite := vc.favoriteCache.RedisIsUserLikeVideosCreated(int64(userId), int64(videoId))
	var isFav bool
	if isFavorite == codes.BITMAPLIKE {
		isFav = true
	} else if isFavorite == codes.BITMAPUNLIKE {
		isFav = false
	} else if isFavorite == codes.ERROR { // 查数据库，如果没查到默认false
		flag, fav := favoriteRepository.GetFavoriteByUserIdAndVideoId(userId, videoId)
		if flag {
			isFav = fav.IsDeleted
		}
		isFav = false
	}
	var favoriteCount uint32
	redisFavCount, _ := vc.favoriteCache.RedisGetVideoFavoriteCount(int64(videoId))
	if redisFavCount == 0 {
		video := videoRepository.GetVideoByVideoId(videoId) // 查数据库
		favoriteCount = video.FavoriteCount
	} else {
		favoriteCount = uint32(redisFavCount)
	}
	return favoriteCount, isFav
}

// getCommentCount 查评论数量
func (vc *VideoCache) getCommentCount(videoId uint64) uint32 {
	CommentString := "videoComment:comment"
	var commentVos []vo.CommentVo
	var commentCount uint32
	data, _ := global.REDIS.Get(context.Background(), CommentString+strconv.FormatUint(videoId, 10)).Result()
	if data != "" {
		err := json.Unmarshal([]byte(data), &commentVos)
		if err != nil {
			log.Println(err)
		}
		commentCount = uint32(len(commentVos))
	} else { // 从数据库中找
		video := videoRepository.GetVideoByVideoId(videoId)
		commentCount = video.CommentCount
	}
	return commentCount
}

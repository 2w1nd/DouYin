package service

import (
	"encoding/json"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
)

type FeedService struct {
}

var videoRepository repository.VideoRepository
var favoriteService FavoriteService
var userService UserService

type VideoMsg struct {
	VideoID    uint64 `json:"id,omitempty"`
	AuthorID   uint64 `json:"author_id,omitempty"`
	PlayUrl    string `json:"play_url,omitempty"`
	CoverUrl   string `json:"cover_url,omitempty"`
	Title      string `json:"title,omitempty"`
	CreateTime int64  `json:"create_time"`
}

// Feed
// @Description: 视频流接口
// @receiver: fs
// @param: token
// @param: latestTime
// @return: []vo.VideoVo
func (fs *FeedService) Feed(userId uint64, latestTime string) ([]vo.VideoVo, time.Time) {
	var (
		//videoData vo.VideoData
		videoVos  []vo.VideoVo
		videoList []model.Video
		nextTime  int64
	)
	//从缓存查询
	videoVos, nextTime = fs.readFeedDataFromRedis(userId)
	if len(videoVos) != 0 {
		return videoVos, time.Unix(nextTime/1000, 0)
	}

	//data1, _ := global.REDIS.Get(context.Background(), "videoVos").Result()
	//if data1 != "" {
	//	log.Println("从缓存中查询: ", data1)
	//	err := json.Unmarshal([]byte(data1), &videoData)
	//	if err != nil {
	//		return nil, time.Time{}
	//	}
	//	if len(videoData.VideoList) != 0 {
	//		return videoData.VideoList, time.Unix(videoData.NextTime/1000, 0)
	//	}
	//}

	// 从数据库中查询
	log.Println("从数据库中查询")
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
	fs.loadFeedDataToRedis(videoList)
	//videoData.VideoList = videoVos
	//videoData.NextTime = timeUtil.TimeToUnix(videoList[0].GmtCreated)
	//data, _ := json.Marshal(videoData)
	//global.REDIS.Set(context.Background(), "videoVos", data, 10*time.Minute)
	return videoVos, videoList[0].GmtCreated
}

func (fs *FeedService) readFeedDataFromRedis(userId uint64) (videoVos []vo.VideoVo, nextTime int64) {
	// 查视频数据
	var videoMsgs []VideoMsg
	videosIds, _ := global.REDIS.LRange(context.Background(), "videoIds", 0, -1).Result()
	for _, videoId := range videosIds {
		video := global.REDIS.HGet(context.Background(), "videos", videoId).String()
		videom := utils.SplitStringForList(video, ":")
		var v VideoMsg
		err := json.Unmarshal([]byte(videom), &v)
		log.Println(video)
		log.Println(videom)
		if err != nil {
			log.Println(err)
		}
		videoMsgs = append(videoMsgs, v)
	}

	for _, videoMsg := range videoMsgs {
		// 查作者名称
		authorName := fs.getAuthorNameInRedis(videoMsg.AuthorID)
		// 查粉丝数量，关注数量，当前用户是否关注
		fs.getFollowCountAndFollowedCountAndIsFollow()
		// 查点赞数量，当前用户是否点赞
		favoriteCount, isFavorite := fs.getFavoriteCountAndIsFavorite(userId, videoMsg.VideoID)
		log.Println("点赞数量：", favoriteCount, videoMsg.VideoID)
		// 查评论数量
		commentCount := fs.getCommentCount(videoMsg.VideoID)
		log.Println("评论数量：", commentCount, videoMsg.VideoID)

		videoVo := vo.VideoVo{
			VideoID: videoMsg.VideoID,
			Author: vo.AuthorVo{
				UserID:        videoMsg.AuthorID,
				Name:          authorName,
				FollowCount:   0,
				FollowerCount: 0,
				IsFollow:      false,
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

// 查作者名称
func (fs *FeedService) getAuthorNameInRedis(authorId uint64) string {
	log.Println("作者名称")
	authorName, _ := global.REDIS.HGet(context.Background(), "users:user", strconv.FormatUint(authorId, 10)).Result()
	log.Println(authorName)
	if authorName == "" { // 找不到，从数据库找
		authorName = userService.GetUserName(authorId)
	}
	return authorName
}

// 查粉丝数量，关注数量，当前用户是否关注
func (fs *FeedService) getFollowCountAndFollowedCountAndIsFollow() {

}

// 查点赞数量，当前用户是否点赞
func (fs *FeedService) getFavoriteCountAndIsFavorite(userId uint64, videoId uint64) (uint32, bool) {
	isFavorite := favoriteService.RedisIsUserLikeVideosCreated(int64(userId), int64(videoId))
	var isFav bool
	if isFavorite == codes.BITMAPLIKE {
		isFav = true
	} else if isFavorite == codes.BITMAPUNLIKE {
		isFav = false
	} else if isFavorite == codes.ERROR { // 查数据库，如果没查到默认false
		isFav = favoriteService.DBIsUserLikeVideosCreated(int64(userId), int64(videoId))
	}
	var favoriteCount uint32
	redisFavCount, _ := RedisGetVideoFavoriteCount(int64(videoId))
	if redisFavCount == 0 {
		video := videoRepository.GetVideoByVideoId(videoId) // 查数据库
		favoriteCount = video.FavoriteCount
	} else {
		favoriteCount = uint32(redisFavCount)
	}
	return favoriteCount, isFav
}

// 查评论数量
func (fs *FeedService) getCommentCount(videoId uint64) uint32 {
	log.Println("评论数量")
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

func (fs *FeedService) loadFeedDataToRedis(videoModels []model.Video) {
	for _, video := range videoModels {
		videomsg := VideoMsg{
			VideoID:    video.VideoId,
			AuthorID:   video.AuthorId,
			PlayUrl:    video.Path,
			CoverUrl:   video.CoverPath,
			Title:      video.Title,
			CreateTime: utils.TimeToUnix(video.GmtCreated),
		}
		//videoIdString := strconv.FormatUint(video.VideoId, 10)
		videoMsgJson, _ := json.Marshal(videomsg)
		//log.Println(videoMsgJson)
		str := global.REDIS.HGet(context.Background(), "videos", strconv.FormatUint(video.VideoId, 10)).Err()
		if str != nil {
			global.REDIS.HMSet(context.Background(), "videos", strconv.FormatUint(video.VideoId, 10), videoMsgJson)
			global.REDIS.RPush(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10))
		}
	}
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

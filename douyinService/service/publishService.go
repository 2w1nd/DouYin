package service

import (
	"context"
	"encoding/json"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/convert"
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"strconv"
	"sync"
)

type PublishService struct {
}

var userRepository repository.UserRepository

func (ps *PublishService) Publish(userId uint64, data *multipart.FileHeader, title string) {
	//上传文件
	newKey := uuid.New().String()
	task := sync.WaitGroup{}
	task.Add(1)
	go func() {
		utils.UploadVideo(newKey, data)
		task.Done()
	}()
	log.Println("上次文件，准备存入数据库")
	path := "http://img.xlong.xyz/video/" + newKey
	cover := path + "?vframe/jpg/offset/1"
	//存入数据库
	log.Println(userId)
	video := model.Video{
		VideoId:       uint64(global.ID.Generate()),
		AuthorId:      userId,
		Title:         title,
		Path:          path,
		CoverPath:     cover,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	videoRepository.SaveVideo(video)
	ps.loadPublishDataToRedis(video)
	task.Wait()
}

func (ps *PublishService) PublishList(myId uint64, userId uint64) []dto.VideoDto {
	//查询目标用户信息,只查一遍
	user := userRepository.QueryUserDtoInfo(userId)
	userDto := convert.User2UserDTO(user)
	var videoDtoList []dto.VideoDto

	// 读缓存
	videoVos := ps.readPublishDataFromRedis(userId, myId)
	if len(videoVos) != 0 {
		return videoVos
	}

	log.Println("查询数据库")
	if myId == 0 {
		videoDtoList = ps.publishListWithoutLogin(userId)
	} else {
		//如果登录了，则填充IsFollow
		userDto.IsFollow = userRepository.IsFollow(myId, userId)
		videoDtoList = ps.publishListWithLogin(myId, userId)
	}

	//返回结果
	//给每个video填充作者信息
	for i := range videoDtoList {
		videoDtoList[i].Author = userDto
	}
	return videoDtoList
}

func (ps *PublishService) publishListWithLogin(loginUser, userId uint64) []dto.VideoDto {
	//获取用户上传的视频列表
	videoList := videoRepository.GetPublishListWithFavorite(userId, 1, 30, loginUser)
	for _, video := range videoList {
		ps.loadPublishDataToRedis(video)
	}
	videoDtoList := convert.VideoList2VideoDtoList(videoList)

	return videoDtoList
}

func (ps *PublishService) publishListWithoutLogin(userId uint64) []dto.VideoDto {
	//获取用户上传的视频列表
	videoList := videoRepository.GetPublishList(userId, 1, 30)
	for _, video := range videoList {
		ps.loadPublishDataToRedis(video)
	}
	videoDtoList := convert.VideoList2VideoDtoList(videoList)

	return videoDtoList
}

func (ps *PublishService) readPublishDataFromRedis(userId, myId uint64) []dto.VideoDto {
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

	for _, videoMsg := range videoMsgs {
		// 查作者名称
		authorName := ps.getAuthorNameInRedis(videoMsg.AuthorID)
		// 查粉丝数量，关注数量，当前用户是否关注
		followerCount, followCount, isFollow := ps.getFollowCountAndFollowedCountAndIsFollow(myId, userId)
		// 查点赞数量，当前用户是否点赞
		favoriteCount, isFavorite := ps.getFavoriteCountAndIsFavorite(myId, videoMsg.VideoID)
		// 查评论数量
		commentCount := ps.getCommentCount(videoMsg.VideoID)

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

// 查作者名称
func (ps *PublishService) getAuthorNameInRedis(authorId uint64) string {
	authorName, _ := global.REDIS.HGet(context.Background(), "users:user", strconv.FormatUint(authorId, 10)).Result()
	if authorName == "" { // 找不到，从数据库找
		authorName = userService.GetUserName(authorId)
	}
	return authorName
}

// 查粉丝数量，关注数量，当前用户是否关注
func (ps *PublishService) getFollowCountAndFollowedCountAndIsFollow(myId uint64, authorId uint64) (uint32, uint32, bool) {
	followerCount, _ := relationService.RedisGetFollowerCount(int64(authorId))
	followCount, _ := relationService.RedisGetFollowCount(int64(authorId))
	isFollow := relationService.RedisIsRelationCreated(int64(myId), int64(authorId))
	return uint32(followerCount), uint32(followCount), isFollow
}

// 查点赞数量，当前用户是否点赞
func (ps *PublishService) getFavoriteCountAndIsFavorite(userId uint64, videoId uint64) (uint32, bool) {
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
func (ps *PublishService) getCommentCount(videoId uint64) uint32 {
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

func (ps *PublishService) loadPublishDataToRedis(video model.Video) {
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
	str := global.REDIS.HGet(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10)).Err()
	if str != nil {
		global.REDIS.RPush(context.Background(), "video", videoMsgJson)
		global.REDIS.RPush(context.Background(), "userVideos:userVideo"+strconv.FormatUint(video.AuthorId, 10), strconv.FormatUint(video.VideoId, 10))
		global.REDIS.HMSet(context.Background(), "videoIds", strconv.FormatUint(video.VideoId, 10), strconv.FormatUint(video.AuthorId, 10))
	}
}

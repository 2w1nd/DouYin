package service

import (
	"context"
	"encoding/json"
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

	ps.readPublishDataFromRedis(userId, myId)

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

	videoDtoList := convert.VideoList2VideoDtoList(videoList)

	return videoDtoList
}

func (ps *PublishService) publishListWithoutLogin(userId uint64) []dto.VideoDto {
	//获取用户上传的视频列表
	videoList := videoRepository.GetPublishList(userId, 1, 30)

	videoDtoList := convert.VideoList2VideoDtoList(videoList)

	return videoDtoList
}

func (ps *PublishService) readPublishDataFromRedis(userId, myId uint64) vo.VideoVo {
	// 查缓存
	// 得到该用户得视频id列表
	videoIdsStr := global.REDIS.LRange(context.Background(), "userVideos:userVideo"+strconv.FormatUint(userId, 10), 0, -1).String()
	// 遍历视频id
	log.Println(videoIdsStr)
	//for _, videoId := range videoIdsStr {
	//
	//}
	return vo.VideoVo{}
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

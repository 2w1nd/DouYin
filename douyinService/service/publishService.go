package service

import (
	"github.com/DouYin/common/convert"
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/cache"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"sync"
)

type PublishService struct {
	videoCache cache.VideoCache
}

var videoRepository repository.VideoRepository
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
	path := "http://img.xlong.xyz/video/" + newKey
	cover := path + "?vframe/jpg/offset/1"
	//存入数据库
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
	ps.videoCache.LoadPublishDataToRedis(video)
	task.Wait()
}

func (ps *PublishService) PublishList(myId uint64, userId uint64) []dto.VideoDto {
	//查询目标用户信息,只查一遍
	user := userRepository.QueryUserDtoInfo(userId)
	userDto := convert.User2UserDTO(user)
	var videoDtoList []dto.VideoDto

	// 读缓存
	videoVos := ps.videoCache.ReadPublishDataFromRedis(userId, myId)
	videoDtos := convert.VideoVos2VideoDto(videoVos)
	if len(videoDtos) != 0 {
		return videoDtos
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
		ps.videoCache.LoadPublishDataToRedis(video)
	}
	videoDtoList := convert.VideoList2VideoDtoList(videoList)

	return videoDtoList
}

func (ps *PublishService) publishListWithoutLogin(userId uint64) []dto.VideoDto {
	//获取用户上传的视频列表
	videoList := videoRepository.GetPublishList(userId, 1, 30)
	for _, video := range videoList {
		ps.videoCache.LoadPublishDataToRedis(video)
	}
	videoDtoList := convert.VideoList2VideoDtoList(videoList)

	return videoDtoList
}

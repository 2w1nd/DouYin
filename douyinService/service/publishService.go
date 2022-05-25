package service

import (
	"github.com/DouYin/common/context"
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"github.com/google/uuid"
	"mime/multipart"
	"sync"
)

type PublishService struct {
}

var userRepository repository.UserRepository

func (ps *PublishService) Publish(userContext context.UserContext, data *multipart.FileHeader, title string) {
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
		UserId:        userContext.Id,
		Title:         title,
		Path:          path,
		CoverPath:     cover,
		FavoriteCount: 0,
		CommentCount:  0,
	}
	videoRepository.SaveVideo(video)
	task.Wait()
}

func (ps *PublishService) PublishList(userContext context.UserContext, userId uint64) []dto.VideoDto {
	//查询目标用户信息
	userDtoInfo := userRepository.QueryUserDtoInfo(userId)
	userDto := dto.UserDto{
		Id:            userDtoInfo.UserId,
		Name:          userDtoInfo.Name,
		FollowCount:   userDtoInfo.FollowCount,
		FollowerCount: userDtoInfo.FollowerCount,
		IsFollow:      false,
	}
	//如果登录了，则填充IsFollow
	if userContext.Id != 0 {
		userDto.IsFollow = userRepository.IsFollow(userContext.Id, userId)
	}

	//获取用户上传的视频列表
	videoList := videoRepository.GetPublishList(userId)
	//返回结果
	var videoDtoList []dto.VideoDto
	for _, video := range videoList {
		videoDtoList = append(videoDtoList, dto.VideoDto{
			//给每个视频DTO填充Author信息
			Author: userDto,
			Id:     video.VideoId,
		})
	}
	//TODO 填充is_favorite
	return videoDtoList
}

package convert

import (
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
)

func User2UserDTO(user model.User) dto.UserDto {
	return dto.UserDto{
		Id:            user.UserId,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}
}

func Video2VideoDto(video model.Video) dto.VideoDto {
	var isFollow, isFavorite bool
	if len(video.User.FollowedUser) != 0 {
		isFollow = !video.User.FollowedUser[0].IsDeleted
	} else {
		isFollow = false
	}
	if len(video.Favorite) != 0 {
		isFavorite = !video.Favorite[0].IsDeleted
	} else {
		isFavorite = false
	}
	return dto.VideoDto{
		Id: video.VideoId,
		Author: dto.UserDto{
			Id:            video.User.UserId,
			Name:          video.User.Username,
			FollowCount:   video.User.FollowCount,
			FollowerCount: video.User.FollowerCount,
			IsFollow:      isFollow,
		},
		PlayURL:       video.Path,
		CoverURL:      video.CoverPath,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    isFavorite,
	}
}

func VideoList2VideoDtoList(videoList []model.Video) []dto.VideoDto {
	var videoDtoList []dto.VideoDto
	for i := range videoList {
		videoDtoList = append(videoDtoList, Video2VideoDto(videoList[i]))
	}
	return videoDtoList
}

func VideoVos2VideoDto(videoVos []vo.VideoVo) []dto.VideoDto {
	var videoDtos []dto.VideoDto
	for _, videoVo := range videoVos {
		videoDto := dto.VideoDto{
			Id: videoVo.VideoID,
			Author: dto.UserDto{
				Id:            videoVo.Author.UserID,
				Name:          videoVo.Author.Name,
				FollowCount:   videoVo.Author.FollowCount,
				FollowerCount: videoVo.Author.FollowerCount,
				IsFollow:      videoVo.Author.IsFollow,
			},
			PlayURL:       videoVo.PlayUrl,
			CoverURL:      videoVo.CoverUrl,
			Title:         videoVo.Title,
			FavoriteCount: videoVo.FavoriteCount,
			CommentCount:  videoVo.CommentCount,
			IsFavorite:    videoVo.IsFavorite,
		}
		videoDtos = append(videoDtos, videoDto)
	}
	return videoDtos
}

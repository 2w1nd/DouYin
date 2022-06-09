package service

import (
	"context"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/cache"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"strconv"
)

var (
	favoriteRepository repository.FavoriteRepository
	//videoRepository repository.VideoRepository
	ctx = context.Background()
)

type FavoriteService struct {
	favoriteCache cache.FavoriteCache
	videoCache    cache.VideoCache
	userCache     cache.UserCache
}

// DBIsUserLikeVideosCreated 查询db中是否已记录该点赞 方向:User->Videos
func (fs *FavoriteService) DBIsUserLikeVideosCreated(userId int64, videoId int64) bool {
	flag, fav := favoriteRepository.GetFavoriteByUserIdAndVideoId(uint64(userId), uint64(videoId))
	if flag {
		return fav.IsDeleted
	}
	return false
}

// RedisAddFavorite 进行点赞后Redis操作
func (fs *FavoriteService) RedisAddFavorite(favoriteInfo model.Favorite) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = fs.favoriteCache.RedisDeleteUserUnLikeVideos(favoriteInfo); ok == codes.ERROR {
		return false
	}
	if ok = fs.favoriteCache.RedisAddUserLikeVideos(favoriteInfo); ok == codes.ERROR {
		return false
	}
	if ok = fs.favoriteCache.RedisAddVideoLikedByUsers(favoriteInfo); ok == codes.ERROR {
		return false
	}
	return true
}

// RedisDeleteFavorite 取消赞后Redis操作
func (fs *FavoriteService) RedisDeleteFavorite(favoriteInfo model.Favorite) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = fs.favoriteCache.RedisAddUserUnLikeVideos(favoriteInfo); ok == codes.ERROR {
		return false
	}
	if ok = fs.favoriteCache.RedisDeleteUserLikeVideos(favoriteInfo); ok == codes.ERROR {
		return false
	}
	if ok = fs.favoriteCache.RedisDeleteVideoLikedByUsers(favoriteInfo); ok == codes.ERROR {
		return false
	}
	return true
}

// GetFavoriteList 根据UserId获取用户点赞列表
func (fs *FavoriteService) GetFavoriteList(userId int64) ([]vo.VideoVo, error) {
	var videoVoList []vo.VideoVo
	videoIds, err := fs.RedisGetFavoriteIDs(userId)
	if err != nil {
		return videoVoList, err
	}
	videoVoList, _ = fs.videoCache.GetVideoVoByIdsFromRedis(videoIds, uint64(userId))
	return videoVoList, nil
}

// RedisGetFavoriteIDs 返回用户喜欢视频IDs
func (fs *FavoriteService) RedisGetFavoriteIDs(userId int64) ([]string, error) {
	var videoIds []string
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctx, zsetKey, 0, -1).Result()

	if err != nil {
		return videoIds, err
	}
	for _, value := range values {
		//Member为interface类型不能进行强制转换
		videoid, _ := value.Member.(string)
		videoIds = append(videoIds, videoid)
	}
	return videoIds, err
}

func (fs *FavoriteService) videoList2Vo(videoList []model.Video) []vo.VideoVo {
	var videoVos []vo.VideoVo
	for _, video := range videoList {
		var isFollow, isFavorite bool
		if len(video.User.FollowedUser) != 0 {
			isFollow = video.User.FollowedUser[0].IsDeleted
		} else {
			isFollow = false
		}
		ok := fs.favoriteCache.RedisIsUserLikeVideosCreated(int64(video.User.UserId), int64(video.VideoId))
		if ok == codes.BITMAPLIKE {
			isFavorite = true
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
		count, _ := fs.favoriteCache.RedisGetVideoFavoriteCount(int64(videoVo.VideoID))
		videoVo.FavoriteCount = uint32(count)

		videoVos = append(videoVos, videoVo)
	}
	return videoVos
}

package service

import (
	"context"
	"fmt"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	//bitmap查询结果
	BITMAPLIKE   int = 1
	BITMAPUNLIKE int = 0

	ERROR int = -1

	ALREADYEXIST  int = 0
	ALREADYDELETE int = 0
	SUCCESS       int = 1
)

var (
	favoriteRepository repository.FavoriteRepository
	ctx                = context.Background()
)

type FavoriteService struct {
}

//Time转为float64类型的时间戳
func Time2Float(time time.Time) float64 {
	return float64(time.Unix())
}

//点赞后，zset添加用户取消点赞的VideoId(score为时间) 方向:User->Videos
func redisAddUserLikeVideos(favoriteInfo model.Favorite) int {
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(favoriteInfo.UserId))

	zsetScore := Time2Float(time.Now())
	zsetMember := strconv.Itoa(int(favoriteInfo.VideoId))
	zsetValue := &redis.Z{zsetScore, zsetMember}

	_, err := global.REDIS.ZRank(ctx, zsetKey, zsetMember).Result()
	//未找到对应的key或未找到对应的key和value,返回err=redis:null
	if err == redis.Nil {
		fmt.Println("error", err)
		_, err = global.REDIS.ZAdd(ctx, zsetKey, zsetValue).Result()

		fmt.Println("redisAddUserLikeVideos", 0)
		if err != nil {
			fmt.Println("redisAddUserLikeVideos", 1)
			return ERROR
		}
		return SUCCESS
	} else if err != nil {
		return ERROR
	} else {
		//key和value已存在

		fmt.Println("redisAddUserLikeVideos", 2)
		return ALREADYEXIST

		return SUCCESS
	}

}

//点赞后取消，set添加用户取消点赞的VideoId 方向:User->Videos
func redisAddUserUnLikeVideos(favoriteInfo model.Favorite) int {
	setKey := "favorite:" + "user_unlike_videos:" + strconv.Itoa(int(favoriteInfo.UserId))
	setValue := strconv.Itoa(int(favoriteInfo.VideoId))

	//未查找到该取消赞的VideoId,key不存在时返回0
	ok, err := global.REDIS.SIsMember(ctx, setKey, setValue).Result()
	//fmt.Println("redisAddUserUnLikeVideos", 1)
	if err != nil {
		return ERROR
		//fmt.Println("redisAddUserUnLikeVideos", 2)
	}
	if ok == true {
		//fmt.Println("redisAddUserUnLikeVideos", 3)
		return ALREADYEXIST
	}
	//fmt.Println("redisAddUserUnLikeVideos", 4)
	//添加取消赞的VideoId
	_, err = global.REDIS.SAdd(ctx, setKey, setValue).Result()
	if err != nil {
		//fmt.Println("redisAddUserUnLikeVideos", 5)
		return ERROR
	}
	//fmt.Println("redisAddUserUnLikeVideos", 6)
	return SUCCESS
}

//查询bitmap中是否已记录该点赞 方向:User->Videos
func (fs *FavoriteService) RedisIsUserLikeVideosCreated(userId int64, videoId int64) int {
	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(videoId))

	ok, err := global.REDIS.GetBit(ctx, bitmapKey, userId).Result()
	//查询失败
	if err != nil {
		fmt.Println("RedisIsUserLikeVideosCreated", 1)
		return ERROR
	}

	//返回查询到的结果
	if ok == 1 {
		fmt.Println("RedisIsUserLikeVideosCreated", 2)
		return BITMAPLIKE
	} else {
		//未记录过或值为0
		fmt.Println("RedisIsUserLikeVideosCreated", 2)
		return BITMAPUNLIKE
	}

}

//取消赞后，记录取消赞 方向:User->Videos
func redisDeleteUserLikeVideos(favoriteInfo model.Favorite) int {
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(favoriteInfo.UserId))

	zsetMember := strconv.Itoa(int(favoriteInfo.VideoId))

	_, err := global.REDIS.ZRank(ctx, zsetKey, zsetMember).Result()
	//点赞结果已删除
	if err == redis.Nil {
		//fmt.Println("redisDeleteUserLikeVideos", 0)
		return ALREADYDELETE
	} else if err != nil {
		return ERROR
	} else {

		//_位置返回1删除成功，0是不存在
		_, err = global.REDIS.ZRem(ctx, zsetKey, zsetMember).Result()
		//fmt.Println("redisDeleteUserLikeVideos", 1)
		if err != nil {
			//fmt.Println("redisDeleteUserLikeVideos", err)
			//fmt.Println("redisDeleteUserLikeVideos", 2)
			return ERROR
		}

		return SUCCESS

	}

}

//取消赞后又点赞时，在记录删除取消赞的set中删除 方向:User->Videos
func redisDeleteUserUnLikeVideos(favoriteInfo model.Favorite) int {
	setKey := "favorite:" + "user_unlike_videos:" + strconv.Itoa(int(favoriteInfo.UserId))
	setValue := strconv.Itoa(int(favoriteInfo.VideoId))

	//未查找到该取消赞的VideoId
	ok, err := global.REDIS.SIsMember(ctx, setKey, setValue).Result()
	if err != nil {
		return ERROR
	}
	if ok == false {
		return ALREADYDELETE
	}

	_, err = global.REDIS.SRem(ctx, setKey, setValue).Result()
	if err != nil {
		return ERROR
	}

	return SUCCESS
}

//点赞后，bitmap将该UserId位置1 方向:Video->Users
func redisAddVideoLikedByUsers(favoriteInfo model.Favorite) int {

	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(favoriteInfo.VideoId))

	_, err := global.REDIS.SetBit(ctx, bitmapKey, int64(favoriteInfo.UserId)%4294967296, 1).Result()
	if err != nil {
		return ERROR
	}
	return SUCCESS

}

//取消赞后，bitmap将该UserId位置0 方向:Video->Users
func redisDeleteVideoLikedByUsers(favoriteInfo model.Favorite) int {

	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(favoriteInfo.VideoId))

	_, err := global.REDIS.SetBit(ctx, bitmapKey, int64(favoriteInfo.UserId)%4294967296, 0).Result()
	if err != nil {
		//fmt.Println("redisDeleteVideoLikedByUsers", 1)
		return ERROR
	}
	//fmt.Println("redisDeleteVideoLikedByUsers", 0)
	return SUCCESS

}

//进行点赞后Redis操作
func (fs *FavoriteService) RedisAddFavorite(favoriteInfo model.Favorite) bool {

	var ok int
	//后面加锁，保证原子性
	if ok = redisDeleteUserUnLikeVideos(favoriteInfo); ok == ERROR {
		return false
	}
	if ok = redisAddUserLikeVideos(favoriteInfo); ok == ERROR {
		return false
	}
	if ok = redisAddVideoLikedByUsers(favoriteInfo); ok == ERROR {
		return false
	}

	return true

}

//取消赞后Redis操作
func (fs *FavoriteService) RedisDeleteFavorite(favoriteInfo model.Favorite) bool {

	var ok int
	//后面加锁，保证原子性
	if ok = redisAddUserUnLikeVideos(favoriteInfo); ok == ERROR {
		return false
	}
	if ok = redisDeleteUserLikeVideos(favoriteInfo); ok == ERROR {
		return false
	}
	if ok = redisDeleteVideoLikedByUsers(favoriteInfo); ok == ERROR {
		return false
	}

	return true

}

func (fs *FavoriteService) RedisGetFavoriteList(userId int64) ([]int64, error) {
	var videoIds []int64
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctx, zsetKey, 0, -1).Result()

	if err != nil {
		return videoIds, err
	}
	//fmt.Println("values:", values)
	for _, value := range values {
		//Member为interface类型不能进行强制转换
		videoid, _ := strconv.ParseInt(value.Member.(string), 10, 64)
		videoIds = append(videoIds, videoid)
	}
	//fmt.Println("VideoId:", videoIds)

	return videoIds, err

}

//根据VideoId查询对应Video点赞数量
func (fs *FavoriteService) RedisGetVideoFavoriteCount(videoId int64) (int64, error) {

	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(videoId))
	var favoriteCount *redis.BitCount
	ans, err := global.REDIS.BitCount(ctx, bitmapKey, favoriteCount).Result()
	return ans, err
}

//根据UserId获取用户点赞列表
func (fs *FavoriteService) GetFavoriteList(userId int64) ([]vo.VideoVo, error) {

	var videoList []model.Video
	var videoVoList []vo.VideoVo

	videoIds, err := fs.RedisGetFavoriteList(userId)
	if err != nil {
		return videoVoList, err
	}
	fmt.Println("videoIds::", videoIds)
	for _, id := range videoIds {
		videoList = append(videoList, videoRepository.GetVideoByVideoId(uint64(id)))
	}
	videoVoList = fs.videoList2Vo(videoList)
	return videoVoList, nil

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
		ok := fs.RedisIsUserLikeVideosCreated(int64(video.User.UserId), int64(video.VideoId))
		if ok == BITMAPLIKE {
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
		count, _ := fs.RedisGetVideoFavoriteCount(int64(videoVo.VideoID))
		videoVo.FavoriteCount = uint32(count)

		videoVos = append(videoVos, videoVo)
	}
	return videoVos
}

func SynchronizeDBAndRedis() {

	zsetkey, err := global.REDIS.Keys(ctx, "favorite:"+"user_like_videos:*").Result()
	if err != nil {
		return
	}
	var FavoritevideoIds []string
	for userId := range zsetkey {
		FavoritevideoIds, err = global.REDIS.ZRange(ctx, strconv.Itoa(userId), 0, -1).Result()
		for videoId := range FavoritevideoIds {
			var ok bool
			ok, _ = favoriteRepository.GetFavoriteByUserIdAndVideoId(int64(userId), int64(videoId))
			if ok {
				return
			}
			favoriteRepository.AddFavorite(model.Favorite{
				UserId:    uint64(userId),
				VideoId:   uint64(videoId),
				IsDeleted: true,
			})
		}

	}

	zsetkey, err = global.REDIS.Keys(ctx, "favorite:"+"user_unlike_videos:*").Result()
	if err != nil {
		return
	}
	var UnFavoritevideoIds []string
	for userId := range zsetkey {
		UnFavoritevideoIds, err = global.REDIS.ZRange(ctx, strconv.Itoa(userId), 0, -1).Result()
		for videoId := range UnFavoritevideoIds {
			var ok bool
			ok, _ = favoriteRepository.GetFavoriteByUserIdAndVideoId(int64(userId), int64(videoId))
			if ok {
				return
			}
			favoriteRepository.AddFavorite(model.Favorite{
				UserId:    uint64(userId),
				VideoId:   uint64(videoId),
				IsDeleted: true,
			})
		}

	}

}

package cache

import (
	"context"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"github.com/go-redis/redis/v8"
	"log"
	"math"
	"strconv"
	"time"
)

var favoriteRepository repository.FavoriteRepository
var ctx = context.Background()

type FavoriteCache struct {
}

// Time2Float Time转为float64类型的时间戳
func Time2Float(time time.Time) float64 {
	return float64(time.Unix())
}

// RedisAddUserLikeVideos 点赞后，zset添加用户取消点赞的VideoId(score为时间) 方向:User->Videos
func (fc *FavoriteCache) RedisAddUserLikeVideos(favoriteInfo model.Favorite) int {
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(favoriteInfo.UserId))

	zsetScore := Time2Float(time.Now())
	zsetMember := strconv.Itoa(int(favoriteInfo.VideoId))
	zsetValue := &redis.Z{zsetScore, zsetMember}

	_, err := global.REDIS.ZRank(ctx, zsetKey, zsetMember).Result()
	//未找到对应的key或未找到对应的key和value,返回err=redis:null
	if err == redis.Nil {
		_, err = global.REDIS.ZAdd(ctx, zsetKey, zsetValue).Result()
		if err != nil {
			return codes.ERROR
		}
		return codes.SUCCESS
	} else if err != nil {
		return codes.ERROR
	} else {
		//key和value已存在
		return codes.ALREADYEXIST
	}
}

// RedisAddUserUnLikeVideos 点赞后取消，set添加用户取消点赞的VideoId 方向:User->Videos
func (fc *FavoriteCache) RedisAddUserUnLikeVideos(favoriteInfo model.Favorite) int {
	setKey := "favorite:" + "user_unlike_videos:" + strconv.Itoa(int(favoriteInfo.UserId))
	setValue := strconv.Itoa(int(favoriteInfo.VideoId))

	//未查找到该取消赞的VideoId,key不存在时返回0
	ok, err := global.REDIS.SIsMember(ctx, setKey, setValue).Result()
	if err != nil {
		return codes.ERROR
	}
	if ok == true {
		return codes.ALREADYEXIST
	}
	//添加取消赞的VideoId
	_, err = global.REDIS.SAdd(ctx, setKey, setValue).Result()
	if err != nil {
		return codes.ERROR
	}
	return codes.SUCCESS
}

// RedisIsUserLikeVideosCreated 查询bitmap中是否已记录该点赞 方向:User->Videos
func (fc *FavoriteCache) RedisIsUserLikeVideosCreated(userId int64, videoId int64) int {
	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(videoId))
	ok, err := global.REDIS.GetBit(ctx, bitmapKey, userId%math.MaxInt32).Result()
	//查询失败
	if err != nil {
		return codes.ERROR
	}
	//返回查询到的结果
	if ok == 1 {
		return codes.BITMAPLIKE
	} else {
		//未记录过或值为0
		return codes.BITMAPUNLIKE
	}
}

// RedisDeleteUserLikeVideos 取消赞后，记录取消赞 方向:User->Videos
func (fc *FavoriteCache) RedisDeleteUserLikeVideos(favoriteInfo model.Favorite) int {
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(favoriteInfo.UserId))

	zsetMember := strconv.Itoa(int(favoriteInfo.VideoId))

	_, err := global.REDIS.ZRank(ctx, zsetKey, zsetMember).Result()
	//点赞结果已删除
	if err == redis.Nil {
		return codes.ALREADYDELETE
	} else if err != nil {
		return codes.ERROR
	} else {
		//_位置返回1删除成功，0是不存在
		_, err = global.REDIS.ZRem(ctx, zsetKey, zsetMember).Result()
		if err != nil {
			return codes.ERROR
		}
		return codes.SUCCESS
	}
}

// RedisDeleteUserUnLikeVideos 取消赞后又点赞时，在记录删除取消赞的set中删除 方向:User->Videos
func (fc *FavoriteCache) RedisDeleteUserUnLikeVideos(favoriteInfo model.Favorite) int {
	setKey := "favorite:" + "user_unlike_videos:" + strconv.Itoa(int(favoriteInfo.UserId))
	setValue := strconv.Itoa(int(favoriteInfo.VideoId))

	//未查找到该取消赞的VideoId
	ok, err := global.REDIS.SIsMember(ctx, setKey, setValue).Result()
	if err != nil {
		return codes.ERROR
	}
	if ok == false {
		return codes.ALREADYDELETE
	}

	_, err = global.REDIS.SRem(ctx, setKey, setValue).Result()
	if err != nil {
		return codes.ERROR
	}

	return codes.SUCCESS
}

// RedisAddVideoLikedByUsers 点赞后，bitmap将该UserId位置1 方向:Video->Users
func (fc *FavoriteCache) RedisAddVideoLikedByUsers(favoriteInfo model.Favorite) int {
	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(favoriteInfo.VideoId))
	_, err := global.REDIS.SetBit(ctx, bitmapKey, int64(favoriteInfo.UserId)%math.MaxInt32, 1).Result()
	if err != nil {
		return codes.ERROR
	}
	return codes.SUCCESS
}

// RedisDeleteVideoLikedByUsers 取消赞后，bitmap将该UserId位置0 方向:Video->Users
func (fc *FavoriteCache) RedisDeleteVideoLikedByUsers(favoriteInfo model.Favorite) int {
	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(favoriteInfo.VideoId))
	_, err := global.REDIS.SetBit(ctx, bitmapKey, int64(favoriteInfo.UserId)%math.MaxInt32, 0).Result()
	if err != nil {
		return codes.ERROR
	}
	return codes.SUCCESS
}

// RedisGetVideoFavoriteCount 根据VideoId查询对应Video点赞数量
func (fc *FavoriteCache) RedisGetVideoFavoriteCount(videoId int64) (uint32, int) {
	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(videoId))
	var favoriteCount *redis.BitCount
	ans, err := global.REDIS.BitCount(ctx, bitmapKey, favoriteCount).Result()
	if err == redis.Nil {
		return 0, codes.RedisNotFound
	}
	return uint32(ans), codes.RedisNotFound
}

// GetFavoriteCountAndIsFavorite 查点赞数量，当前用户是否点赞
func (fc *FavoriteCache) GetFavoriteCountAndIsFavorite(userId uint64, videoId uint64) (uint32, bool) {
	isFavorite := fc.RedisIsUserLikeVideosCreated(int64(userId), int64(videoId))
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
	favoriteCount, _ = fc.RedisGetVideoFavoriteCount(int64(videoId))
	//if code == codes.RedisNotFound {
	//	log.Println("favoritecache：GetFavoriteCountAndIsFavorite查数据库")
	//	video := videoRepository.GetVideoByVideoId(videoId) // 查数据库
	//	favoriteCount = video.FavoriteCount
	//} else {
	//	favoriteCount = uint32(redisFavCount)
	//}
	return favoriteCount, isFav
}

func SynchronizeFavoriteDBFromRedis() {
	log.Println("同步redis点赞信息到数据库")
	zsetkey, err := global.REDIS.Keys(ctx, "favorite:"+"user_like_videos:*").Result()
	if err != nil {
		return
	}
	var FavoriteVideoIds []string
	for _, userId := range zsetkey {
		FavoriteVideoIds, err = global.REDIS.ZRange(ctx, userId, 0, -1).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		for _, videoId := range FavoriteVideoIds {
			vid := utils.String2Uint64(videoId)
			favoriteRepository.AddFavorite(model.Favorite{
				UserId:    uid,
				VideoId:   vid,
				IsDeleted: false,
			})
		}
	}

	var setkey []string
	setkey, err = global.REDIS.Keys(ctx, "favorite:"+"user_unlike_videos:*").Result()
	if err != nil {
		return
	}
	var UnFavoritevideoIds []string
	for _, userId := range setkey {
		UnFavoritevideoIds, err = global.REDIS.SMembers(ctx, userId).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		for _, videoId := range UnFavoritevideoIds {
			vid := utils.String2Uint64(videoId)
			favoriteRepository.AddFavorite(model.Favorite{
				UserId:    uid,
				VideoId:   vid,
				IsDeleted: true,
			})
		}

	}
	var fc FavoriteCache
	var bitmapkey []string
	bitmapkey, err = global.REDIS.Keys(ctx, "favorite:"+"video_likedby_users:*").Result()
	for _, videoId := range bitmapkey {
		vid := utils.String2Uint64(utils.SplitString(videoId, ":"))
		ans, _ := fc.RedisGetVideoFavoriteCount(int64(vid))
		err := favoriteRepository.SaveFavoriteCount(vid, ans)
		if err != nil {
			return
		}
	}
}

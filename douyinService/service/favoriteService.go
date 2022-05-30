package service

import (
	"context"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/go-redis/redis/v8"
	"log"
	"reflect"
	"strconv"
	"time"
)

const (
	LIKE   bool = false
	UNLIKE bool = true
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

//记录用户点赞的视频id——zset(score为时间)
func redisAddUserLikeVideos(favoriteInfo model.Favorite) *redis.IntCmd {
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(favoriteInfo.UserId))

	zsetScore := Time2Float(time.Now())
	zsetMember := strconv.Itoa(int(favoriteInfo.VideoId))
	zsetValue := &redis.Z{zsetScore, zsetMember}

	var ok *redis.IntCmd
	//为false时记录为点赞
	log.Println(favoriteInfo.IsDeleted)
	if favoriteInfo.IsDeleted == LIKE {
		//未查找到该点赞结果
		if global.REDIS.ZRank(ctx, zsetKey, zsetMember) != nil {
			ok = global.REDIS.ZAdd(ctx, zsetKey, zsetValue)
			log.Println(zsetKey, "点赞成功")
		}
	} else if favoriteInfo.IsDeleted == UNLIKE {
		if global.REDIS.ZRank(ctx, zsetKey, zsetMember) != nil {
			ok = global.REDIS.ZRem(ctx, zsetKey, zsetValue)
			log.Println(zsetKey, "取消赞成功")
		}
	}
	return ok
}

//记录视频对应的点赞用户——bitmap
func redisAddVideLikedByUsers(favoriteInfo model.Favorite) *redis.IntCmd {

	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(favoriteInfo.UserId))

	var ok *redis.IntCmd
	if favoriteInfo.IsDeleted == LIKE {
		ok = global.REDIS.SetBit(ctx, bitmapKey, int64(favoriteInfo.UserId)%10, 1)

	} else {
		ok = global.REDIS.SetBit(ctx, bitmapKey, int64(favoriteInfo.UserId)%10, 0)

	}

	return ok
}

func (fs *FavoriteService) RedisAddFavorite(favoriteInfo model.Favorite) bool {
	var ok *redis.IntCmd
	//后面加锁，保证原子性
	if ok = redisAddUserLikeVideos(favoriteInfo); ok.Err() != nil {
		return false
	}
	if ok = redisAddVideLikedByUsers(favoriteInfo); ok.Err() != nil {
		return false
	}
	return true

}

func (fs *FavoriteService) RedisGetFavoriteList(userId int64) ([]int64, error) {
	var zsetValues []int64
	zsetKey := "favorite:" + "user_like_videos:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctx, zsetKey, 0, -1).Result()

	if err != nil {
		return zsetValues, err
	}
	if reflect.TypeOf(values).Kind() == reflect.String {
		for _, value := range values {
			var num int
			num, err = strconv.Atoi(reflect.ValueOf(value).String())
			zsetValues = append(zsetValues, int64(num))
		}
	}
	return zsetValues, err

}

func (fs *FavoriteService) RedisGetVideoFavoriteCount(videoId int64) (int64, error) {

	bitmapKey := "favorite:" + "video_likedby_users:" + strconv.Itoa(int(videoId))
	var favoriteCount *redis.BitCount
	ans, err := global.REDIS.BitCount(ctx, bitmapKey, favoriteCount).Result()
	return ans, err
}

func (fs *FavoriteService) GetFavoriteList(videoIds []int64) {

}

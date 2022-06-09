package cache

import (
	"context"
	"fmt"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

type RelationCache struct {
}

var (
	followRepository repository.FollowRepository
	ctxx             = context.Background()
)

// RedisGetFollowerCount 根据UserId查询对应用户粉丝数量
func (rc *RelationCache) RedisGetFollowerCount(userid int64) (uint32, int) {
	zsetKey0 := "relation:" + "follower:" + strconv.Itoa(int(userid))
	count, err := global.REDIS.ZCard(ctxx, zsetKey0).Result()
	if err == redis.Nil {
		return 0, codes.RedisNotFound
	}
	return uint32(count), codes.RedisFound
}

// RedisGetFollowCount 根据UserId查询对应用户关注数量
func (rc *RelationCache) RedisGetFollowCount(userid int64) (uint32, int) {
	zsetKey0 := "relation:" + "follow:" + strconv.Itoa(int(userid))
	count, err := global.REDIS.ZCard(ctxx, zsetKey0).Result()
	if err == redis.Nil {
		return 0, codes.RedisNotFound
	}
	return uint32(count), codes.RedisFound
}

// RedisIsRelationCreated 查询中是否已记录该关注 方向:User->ToFollowed
func (rc *RelationCache) RedisIsRelationCreated(userId int64, followedUserId int64) (bool, int) {
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(userId))

	values, err := global.REDIS.ZRevRangeWithScores(ctxx, zsetKey, 0, -1).Result()
	if err != nil {
		return false, codes.RedisNotFound
	}
	for _, value := range values {
		followId, _ := strconv.ParseInt(value.Member.(string), 10, 64)
		if followId == followedUserId {
			return true, codes.RedisFound
		}
	}
	return false, codes.RedisFound
}

// RedisDeleteUserUnRelation 取消关注后又关注时，在记录删除取消关注的set中删除 方向User->ToFollowed
func (rc *RelationCache) RedisDeleteUserUnRelation(followInfo model.Follow) int {
	setKey := "relation:" + "unfollow:" + strconv.Itoa(int(followInfo.UserId))
	setKey0 := "relation:" + "unfollower:" + strconv.Itoa(int(followInfo.FollowedUserId))

	setValue := strconv.Itoa(int(followInfo.FollowedUserId))
	setValue0 := strconv.Itoa(int(followInfo.UserId))

	//未查找到该取消关注的FollowedUserId
	ok, err := global.REDIS.SIsMember(ctxx, setKey, setValue).Result()
	_, _ = global.REDIS.SIsMember(ctxx, setKey0, setValue0).Result()
	if err != nil {
		return codes.ERROR
	}
	if ok == false {
		return codes.ALREADYDELETE
	}

	_, err = global.REDIS.SRem(ctxx, setKey, setValue).Result()
	_, _ = global.REDIS.SRem(ctxx, setKey0, setValue0).Result()
	if err != nil {
		return codes.ERROR
	}

	return codes.SUCCESS
}

// RedisAddRelation 关注后，zset添加用户关注的FollowedUserId(score为时间) 方向:User->ToFollowed
func (rc *RelationCache) RedisAddRelation(relationInfo model.Follow) int {
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(relationInfo.UserId))
	zsetKey0 := "relation:" + "follower:" + strconv.Itoa(int(relationInfo.FollowedUserId))

	zsetScore := Time2Float(time.Now())
	zsetScore0 := zsetScore

	zsetMember := strconv.Itoa(int(relationInfo.FollowedUserId))
	zsetMember0 := strconv.Itoa(int(relationInfo.UserId))

	zsetValue := &redis.Z{Score: zsetScore, Member: zsetMember}
	zsetValue0 := &redis.Z{Score: zsetScore0, Member: zsetMember0}

	_, err := global.REDIS.ZRank(ctxx, zsetKey, zsetMember).Result()
	_, _ = global.REDIS.ZRank(ctxx, zsetKey0, zsetMember0).Result()

	if err == redis.Nil {
		fmt.Println("error", err)
		_, err = global.REDIS.ZAdd(ctxx, zsetKey, zsetValue).Result()
		_, _ = global.REDIS.ZAdd(ctxx, zsetKey0, zsetValue0).Result()

		if err != nil {
			return codes.ERROR
		}
		return codes.SUCCESS
	} else if err != nil {
		return codes.ERROR
	} else {
		// key和value已存在
		return codes.ALREADYEXIST
	}
}

// RedisAddUserUnRelations 关注后取消，set添加用户取消关注的FollowedUserId 方向:User->ToFollowed
func (rc *RelationCache) RedisAddUserUnRelations(relationInfo model.Follow) int {
	setKey := "relation:" + "unfollow:" + strconv.Itoa(int(relationInfo.UserId))
	setKey0 := "relation:" + "unfollower:" + strconv.Itoa(int(relationInfo.FollowedUserId))

	setValue := strconv.Itoa(int(relationInfo.FollowedUserId))
	setValue0 := strconv.Itoa(int(relationInfo.UserId))

	ok, err := global.REDIS.SIsMember(ctxx, setKey, setValue).Result()
	_, _ = global.REDIS.SIsMember(ctxx, setKey0, setValue0).Result()

	if err != nil {
		return codes.ERROR
	}
	if ok == true {
		return codes.ALREADYEXIST
	}

	_, err = global.REDIS.SAdd(ctxx, setKey, setValue).Result()
	_, _ = global.REDIS.SAdd(ctxx, setKey0, setValue0).Result()
	if err != nil {
		return codes.ERROR
	}
	return codes.SUCCESS
}

// RedisUnAddRelation 取消关注后，记录取消的关注 方向:User->ToFollowed
func (rc *RelationCache) RedisUnAddRelation(followInfo model.Follow) int {
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(followInfo.UserId))
	zsetKey0 := "relation:" + "follower:" + strconv.Itoa(int(followInfo.FollowedUserId))

	zsetMember := strconv.Itoa(int(followInfo.FollowedUserId))
	zsetMember0 := strconv.Itoa(int(followInfo.UserId))

	_, err := global.REDIS.ZRank(ctxx, zsetKey, zsetMember).Result()
	_, _ = global.REDIS.ZRank(ctxx, zsetKey0, zsetMember0).Result()
	// 点赞结果已删除
	if err == redis.Nil {
		return codes.ALREADYDELETE
	} else if err != nil {
		return codes.ERROR
	} else {
		//_位置返回1删除成功，0是不存在
		_, err = global.REDIS.ZRem(ctxx, zsetKey, zsetMember).Result()
		_, _ = global.REDIS.ZRem(ctxx, zsetKey0, zsetMember0).Result()
		if err != nil {
			return codes.ERROR
		}
		return codes.SUCCESS
	}
}

// RedisGetFollowList 从redis中获取用户关注用户IDs
func (rc *RelationCache) RedisGetFollowList(userId int64) ([]string, error) {
	var followIds []string
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctxx, zsetKey, 0, -1).Result()

	if err != nil {
		return followIds, err
	}
	for _, value := range values {
		// Member为interface类型不能进行强制转换
		followId, _ := value.Member.(string)
		followIds = append(followIds, followId)
	}
	return followIds, err
}

// RedisGetFollowerList 从redis中获取用户粉丝IDs
func (rc *RelationCache) RedisGetFollowerList(userId int64) ([]string, error) {
	var followerIds []string
	zsetKey := "relation:" + "follower:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctxx, zsetKey, 0, -1).Result()

	if err != nil {
		return followerIds, err
	}
	for _, value := range values {
		// Member为interface类型不能进行强制转换
		followerId, _ := value.Member.(string)
		followerIds = append(followerIds, followerId)
	}
	return followerIds, err
}

// AddAction 添加关注到DB
func (rc *RelationCache) AddAction(where model.Follow) bool {
	var out model.Follow
	followCount, _ := rc.RedisGetFollowCount(int64(where.UserId))
	followerCount, _ := rc.RedisGetFollowerCount(int64(where.UserId))
	if isOk := followRepository.UpdateFollowUserId(where, &out); !isOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowCount(where.UserId, followCount); !IsOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowerCount(where.FollowedUserId, followerCount); !IsOk {
		return false
	}
	return true
}

// RelationAction 取消关注到DB
func (rc *RelationCache) RelationAction(where model.Follow) bool {
	if isOk := followRepository.DeleteFollowUserId(where); !isOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowCount(where.UserId, codes.NoFOCUS); !IsOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowerCount(where.FollowedUserId, codes.NoFOCUS); !IsOk {
		return false
	}
	return true
}

// SynchronizeRelationToDBFromRedis 将redis数据同步到DB
// 分析：对于每个用户，都会有一个Follower zset和一个Follow zset，如当用户A关注用户B时，必然会导致两个zset的改变（B Follower Zset ++，A Follow Zset ++）
// 考虑这种对等情况，所以只需要遍历 Follower key 或 Follow key其中一种入库即可
// 另外还有取关的情况，同样由Unfollower zset 和 Follow zset进行处理
func SynchronizeRelationToDBFromRedis() {
	log.Println("同步redis中关注粉丝信息到数据库")
	zsetkey, err := global.REDIS.Keys(ctxx, "relation:"+"follower:*").Result()
	if err != nil {
		return
	}
	var rs RelationCache
	var FollowerUserIds []string
	for _, userId := range zsetkey {
		FollowerUserIds, err = global.REDIS.ZRange(ctxx, userId, 0, -1).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		for _, FollowedUserId := range FollowerUserIds {
			fid := utils.String2Uint64(FollowedUserId)
			rs.AddAction(model.Follow{
				FollowedUserId: uid,
				UserId:         fid,
				IsDeleted:      false,
			})
		}
	}

	var setkey []string
	setkey, err = global.REDIS.Keys(ctxx, "relation:"+"unfollower:*").Result()
	if err != nil {
		return
	}
	var DeleteFollowedUserIds []string
	for _, userId := range setkey {
		DeleteFollowedUserIds, err = global.REDIS.SMembers(ctxx, userId).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		for _, followedUserId := range DeleteFollowedUserIds {
			vid := utils.String2Uint64(followedUserId)
			rs.RelationAction(model.Follow{
				UserId:         vid,
				FollowedUserId: uid,
			})
		}
	}
}

package service

import (
	"context"
	"fmt"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"github.com/DouYin/service/utils"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
	"time"
)

var (
	followRepository repository.FollowRepository
	ctxx             = context.Background()
)

type RelationService struct {
}

func (rs *RelationService) RelationAction(where model.Follow) bool {
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

// AddAction 添加关注到DB
func (rs *RelationService) AddAction(where model.Follow) bool {
	var out model.Follow
	if isOk := followRepository.UpdateFollowUserId(where, &out); !isOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowCount(where.UserId, codes.FOCUS); !IsOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowerCount(where.FollowedUserId, codes.FOCUS); !IsOk {
		return false
	}

	return true
}

func (rs *RelationService) FollowList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followUsers, _ := followRepository.GetFollowedOrFollowUserWithUserId(userId, codes.Follow)
	log.Println(followUsers)
	for _, user := range followUsers {
		log.Println(user.FollowedUserId)
	}
	userList = rs.userDto2UserVos(followUsers, codes.Follow)
	return userList
}

func (rs *RelationService) FollowerList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followedUsers, _ := followRepository.GetFollowedOrFollowUserWithUserId(userId, codes.Followed)
	log.Println(followedUsers)
	for _, user := range followedUsers {
		log.Println(user.UserId)
	}
	userList = rs.userDto2UserVos(followedUsers, codes.Followed)
	return userList
}

func (rs *RelationService) userDto2UserVos(followerUsers []dto.FollowDto, Type int) []vo.UserVo {
	var userVos []vo.UserVo
	if Type == codes.Followed {
		for _, user := range followerUsers {
			var userVo vo.UserVo
			userVo = vo.UserVo{
				Id:            user.UserId,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      user.FollowedB,
			}
			userVos = append(userVos, userVo)
		}
	} else if Type == codes.Follow {
		for _, user := range followerUsers {
			var userVo vo.UserVo
			userVo = vo.UserVo{
				Id:            user.FollowedUserId,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      user.FollowedB,
			}
			userVos = append(userVos, userVo)
		}
	}
	return userVos
}

// RedisGetFollowerCount 根据UserId查询对应用户粉丝数量
func (rs *RelationService) RedisGetFollowerCount(userid int64) (int64, error) {
	zsetKey0 := "relation:" + "follower:" + strconv.Itoa(int(userid))
	count, err := global.REDIS.ZCard(ctxx, zsetKey0).Result()
	return count, err
}

// RedisGetFollowCount 根据UserId查询对应用户关注数量
func (rs *RelationService) RedisGetFollowCount(userid int64) (int64, error) {
	zsetKey0 := "relation:" + "follow:" + strconv.Itoa(int(userid))
	count, err := global.REDIS.ZCard(ctxx, zsetKey0).Result()
	return count, err
}

// RedisIsRelationCreated 查询中是否已记录该关注 方向:User->ToFollowed
func (rs *RelationService) RedisIsRelationCreated(userId int64, followedUserId int64) bool {
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(userId))

	values, err := global.REDIS.ZRevRangeWithScores(ctxx, zsetKey, 0, -1).Result()
	if err != nil {
		return false
	}
	for _, value := range values {
		followid, _ := strconv.ParseInt(value.Member.(string), 10, 64)
		if followid == followedUserId {
			return true
		}
	}
	return false
}

//取消关注后又关注时，在记录删除取消关注的set中删除 方向User->ToFollowed
func redisDeleteUserUnRelation(followInfo model.Follow) int {
	setKey := "relation:" + "unfollower:" + strconv.Itoa(int(followInfo.UserId))
	setKey0 := "relation:" + "unfollow:" + strconv.Itoa(int(followInfo.FollowedUserId))

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

//关注后，zset添加用户关注的FollowedUserId(score为时间) 方向:User->ToFollowed
func redisAddRelation(relationInfo model.Follow) int {
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
		//key和value已存在
		return codes.ALREADYEXIST
	}
}

// RedisAddRelation 关注后Redis操作
func (rs *RelationService) RedisAddRelation(followInfo model.Follow) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = redisDeleteUserUnRelation(followInfo); ok == codes.ERROR {
		return false
	}
	if ok = redisAddRelation(followInfo); ok == codes.ERROR {
		return false
	}
	return true
}

// 关注后取消，set添加用户取消关注的FollowedUserId 方向:User->ToFollowed
func redisAddUserUnRelations(relationInfo model.Follow) int {
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

// 取消关注后，记录取消的关注 方向:User->ToFollowed
func redisUnAddRelation(followInfo model.Follow) int {
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(followInfo.UserId))
	zsetKey0 := "relation:" + "follower:" + strconv.Itoa(int(followInfo.FollowedUserId))

	zsetMember := strconv.Itoa(int(followInfo.FollowedUserId))
	zsetMember0 := strconv.Itoa(int(followInfo.UserId))

	_, err := global.REDIS.ZRank(ctxx, zsetKey, zsetMember).Result()
	_, _ = global.REDIS.ZRank(ctxx, zsetKey0, zsetMember0).Result()
	//点赞结果已删除
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

// RedisDeleteRelation 取消关注后Redis操作
func (rs *RelationService) RedisDeleteRelation(followInfo model.Follow) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = redisAddUserUnRelations(followInfo); ok == codes.ERROR {
		return false
	}
	if ok = redisUnAddRelation(followInfo); ok == codes.ERROR {
		return false
	}
	return true
}

// RedisGetFollowList 从redis种获取用户关注列表
func (rs *RelationService) RedisGetFollowList(userId int64) ([]int64, error) {
	var followIds []int64
	zsetKey := "relation:" + "followed:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctxx, zsetKey, 0, -1).Result()

	if err != nil {
		return followIds, err
	}
	for _, value := range values {
		//Member为interface类型不能进行强制转换
		followid, _ := strconv.ParseInt(value.Member.(string), 10, 64)
		followIds = append(followIds, followid)
	}
	return followIds, err
}

// RedisGetFollowerList 从redis中获取用户粉丝列表
func (rs *RelationService) RedisGetFollowerList(userId int64) ([]int64, error) {
	var followIds []int64
	zsetKey := "relation:" + "follow:" + strconv.Itoa(int(userId))
	values, err := global.REDIS.ZRevRangeWithScores(ctxx, zsetKey, 0, -1).Result()

	if err != nil {
		return followIds, err
	}
	//fmt.Println("values:", values)
	for _, value := range values {
		//Member为interface类型不能进行强制转换
		followid, _ := strconv.ParseInt(value.Member.(string), 10, 64)
		followIds = append(followIds, followid)
	}

	return followIds, err

}

// GetFollowList 根据UserId获取用户关注列表
func (rs *RelationService) GetFollowList(userId int64) ([]vo.UserVo, error) {
	var followList []model.User
	var followVoList []vo.UserVo

	followIds, err := rs.RedisGetFollowList(userId)
	if err != nil {
		return followVoList, err
	}
	for _, id := range followIds {
		followList = append(followList, userRepository.QueryUserDtoInfo(uint64(id)))
	}
	followVoList = rs.FollowList2Vo(userId, followList)
	return followVoList, nil

}

// GetFollowerList 根据UserId获取用户粉丝列表
func (rs *RelationService) GetFollowerList(userId int64) ([]vo.UserVo, error) {

	var followList []model.User
	var followVoList []vo.UserVo

	followIds, err := rs.RedisGetFollowerList(userId)
	if err != nil {
		return followVoList, err
	}
	for _, id := range followIds {
		followList = append(followList, userRepository.QueryUserDtoInfo(uint64(id)))
	}
	followVoList = rs.FollowList2Vo(userId, followList)
	return followVoList, nil

}

func (rs *RelationService) FollowList2Vo(userId int64, FollowList []model.User) []vo.UserVo {
	var userVos []vo.UserVo
	for _, user := range FollowList {
		var isDelete bool

		ok := rs.RedisIsRelationCreated(int64(user.UserId), userId)
		if ok == true {
			isDelete = true
		} else {
			isDelete = false
		}

		userVo := vo.UserVo{
			Id:       user.UserId,
			Name:     user.Name,
			IsFollow: isDelete,
		}
		count, _ := rs.RedisGetFollowCount(int64(userVo.Id))
		userVo.FollowCount = uint32(count)
		count1, _ := rs.RedisGetFollowerCount(int64(userVo.Id))
		userVo.FollowerCount = uint32(count1)

		userVos = append(userVos, userVo)
	}
	return userVos
}

// SynchronizeRelationToDBFromRedis 将redis数据同步到DB
// 分析：对于每个用户，都会有一个Follower zset和一个Follow zset，如当用户A关注用户B时，必然会导致两个zset的改变（B Follower Zset ++，A Follow Zset ++）
// 考虑这种对等情况，所以只需要遍历 Follower key 或 Follow key其中一种入库即可
// 另外还有取关的情况，同样由Unfollower zset 和 Follow zset进行处理
func SynchronizeRelationToDBFromRedis() {
	log.Println("同步redis到数据库")
	zsetkey, err := global.REDIS.Keys(ctxx, "relation:"+"follower:*").Result()
	log.Println(zsetkey)
	if err != nil {
		return
	}
	var FollowerUserIds []string
	var rs RelationService
	for _, userId := range zsetkey {
		FollowerUserIds, err = global.REDIS.ZRange(ctxx, userId, 0, -1).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		log.Println("like", uid)
		for _, FollowedUserId := range FollowerUserIds {
			log.Println("Followed UserId", FollowedUserId)
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
	log.Println("unfollow setkey: ", setkey)
	var DeleteFollowedUserIds []string
	for _, userId := range setkey {
		DeleteFollowedUserIds, err = global.REDIS.SMembers(ctxx, userId).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		log.Println("unlike:", uid)
		for _, followedUserId := range DeleteFollowedUserIds {
			log.Println("followed userId:", followedUserId)
			vid := utils.String2Uint64(followedUserId)
			rs.RelationAction(model.Follow{
				UserId:         vid,
				FollowedUserId: uid,
			})
		}
	}
}

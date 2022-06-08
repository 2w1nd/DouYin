package service

import (
	"context"
	"github.com/DouYin/common/constant"
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

	ctxx = context.Background()
)

type RelationService struct {
}

const (
	//bitmap查询结果
	BITMAPFOLLOW   int = 1
	BITMAPUNFOLLOW int = 0
)

func (rs *RelationService) RelationAction(where model.Follow) bool {
	if isOk := followRepository.DeleteFollowUserId(where); !isOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowCount(where.UserId, constant.NoFOCUS); !IsOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowerCount(where.FollowedUserId, constant.NoFOCUS); !IsOk {
		return false
	}
	return true
}

func (rs *RelationService) AddAction(where model.Follow) bool {

	var out model.Follow
	if isOk := followRepository.UpdateFollowUserId(where, &out); !isOk {
		follow := model.Follow{
			UserId:         where.UserId,
			FollowedUserId: where.FollowedUserId,
			IsDeleted:      false,
		}
		if IsOk := followRepository.AddFollow(follow); !IsOk {
			return false
		}
	}

	if IsOk := userRepository.UpdateFollowCount(where.UserId, constant.FOCUS); !IsOk {
		return false
	}
	if IsOk := userRepository.UpdateFollowerCount(where.FollowedUserId, constant.FOCUS); !IsOk {
		return false
	}

	return true
}

func (rs *RelationService) FollowList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followUsers, _ := followRepository.GetFollowedOrFollowUserWithUserId(userId, constant.Follow)
	log.Println(followUsers)
	for _, user := range followUsers {
		log.Println(user.FollowedUserId)
	}
	userList = rs.userDto2UserVos(followUsers, constant.Follow)
	return userList
}

func (rs *RelationService) FollowerList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followedUsers, _ := followRepository.GetFollowedOrFollowUserWithUserId(userId, constant.Followed)
	log.Println(followedUsers)
	for _, user := range followedUsers {
		log.Println(user.UserId)
	}
	userList = rs.userDto2UserVos(followedUsers, constant.Followed)
	return userList
}

func (rs *RelationService) userDto2UserVos(followerUsers []dto.FollowDto, Type int) []vo.UserVo {
	var userVos []vo.UserVo
	if Type == constant.Followed {
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
	} else if Type == constant.Follow {
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

//根据UserId查询对应用户粉丝数量
func RedisGetFollowCount(userid int64) (int64, error) {

	bitmapKey := "relation:" + "followedby_users:" + strconv.Itoa(int(userid))
	var relationCount *redis.BitCount
	ans, err := global.REDIS.BitCount(ctxx, bitmapKey, relationCount).Result()
	return ans, err
}

//根据UserId查询对应用户关注数量
func RedisGetFollowedCount(userid int64) (int64, error) {

	bitmapKey := "relation:" + "followby_users:" + strconv.Itoa(int(userid))
	var relationCount *redis.BitCount
	ans, err := global.REDIS.BitCount(ctxx, bitmapKey, relationCount).Result()
	return ans, err
}

//查询bitmap中是否已记录该关注 方向:User->ToFollowed
func (rs *RelationService) RedisIsRelationCreated(userId int64, FollowedUserId int64) int {
	bitmapKey := "relation:" + "followedby_users:" + strconv.Itoa(int(FollowedUserId))
	bitmapKey0 := "relation:" + "followby_users:" + strconv.Itoa(int(userId))

	ok, err := global.REDIS.GetBit(ctxx, bitmapKey, userId).Result()
	global.REDIS.GetBit(ctxx, bitmapKey0, FollowedUserId).Result()
	//查询失败
	if err != nil {
		return ERROR
	}

	//返回查询到的结果
	if ok == 1 {

		return BITMAPFOLLOW
	} else {
		//未记录过或值为0

		return BITMAPUNFOLLOW
	}

}

//取消关注后，记录取消的关注 方向:User->ToFollowed
func redisUnAddRelation(followInfo model.Follow) int {
	zsetKey := "relation:" + "followed:" + strconv.Itoa(int(followInfo.UserId))
	zsetKey0 := "relation:" + "follow:" + strconv.Itoa(int(followInfo.FollowedUserId))

	zsetMember := strconv.Itoa(int(followInfo.FollowedUserId))
	zsetMember0 := strconv.Itoa(int(followInfo.UserId))

	_, err := global.REDIS.ZRank(ctxx, zsetKey, zsetMember).Result()
	global.REDIS.ZRank(ctxx, zsetKey0, zsetMember0).Result()
	//点赞结果已删除
	if err == redis.Nil {
		return ALREADYDELETE
	} else if err != nil {
		return ERROR
	} else {

		//_位置返回1删除成功，0是不存在
		_, err = global.REDIS.ZRem(ctxx, zsetKey, zsetMember).Result()
		global.REDIS.ZRem(ctxx, zsetKey0, zsetMember0).Result()
		if err != nil {
			return ERROR
		}

		return SUCCESS
	}

}

//取消关注后又关注时，在记录删除取消关注的set中删除 方向User->ToFollowed
func redisDeleteUserUnRelation(followInfo model.Follow) int {
	setKey := "relation:" + "unfollowed:" + strconv.Itoa(int(followInfo.UserId))
	setKey0 := "relation:" + "unfollow:" + strconv.Itoa(int(followInfo.FollowedUserId))

	setValue := strconv.Itoa(int(followInfo.FollowedUserId))
	setValue0 := strconv.Itoa(int(followInfo.UserId))

	//未查找到该取消关注的FollowedUserId
	ok, err := global.REDIS.SIsMember(ctxx, setKey, setValue).Result()
	global.REDIS.SIsMember(ctxx, setKey0, setValue0).Result()
	if err != nil {
		return ERROR
	}
	if ok == false {
		return ALREADYDELETE
	}

	_, err = global.REDIS.SRem(ctxx, setKey, setValue).Result()
	global.REDIS.SRem(ctxx, setKey0, setValue0).Result()
	if err != nil {
		return ERROR
	}

	return SUCCESS
}

//关注后，zset添加用户关注的FollowedUserId(score为时间) 方向:User->ToFollowed
func redisAddRelation(relationInfo model.Follow) int {
	zsetKey := "relation:" + "followed:" + strconv.Itoa(int(relationInfo.UserId))
	zsetKey0 := "relation:" + "follow:" + strconv.Itoa(int(relationInfo.FollowedUserId))

	zsetScore := Time2Float(time.Now())
	zsetScore0 := zsetScore

	zsetMember := strconv.Itoa(int(relationInfo.FollowedUserId))
	zsetMember0 := strconv.Itoa(int(relationInfo.UserId))

	zsetValue := &redis.Z{zsetScore, zsetMember}
	zsetValue0 := &redis.Z{zsetScore0, zsetMember0}

	_, err := global.REDIS.ZRank(ctxx, zsetKey, zsetMember).Result()
	global.REDIS.ZRank(ctxx, zsetKey0, zsetMember0).Result()

	if err == redis.Nil {
		//fmt.Println("error", err)
		_, err = global.REDIS.ZAdd(ctxx, zsetKey, zsetValue).Result()
		global.REDIS.ZAdd(ctxx, zsetKey0, zsetValue0).Result()

		if err != nil {
			return ERROR
		}
		return SUCCESS
	} else if err != nil {
		return ERROR
	} else {
		//key和value已存在

		return ALREADYEXIST

	}

}

//关注后取消，set添加用户取消关注的FollowedUserId 方向:User->ToFollowed
func redisAddUserUnRelations(relationInfo model.Follow) int {
	setKey := "relation:" + "unfollowed:" + strconv.Itoa(int(relationInfo.UserId))
	setKey0 := "relation:" + "unfollow:" + strconv.Itoa(int(relationInfo.FollowedUserId))

	setValue := strconv.Itoa(int(relationInfo.FollowedUserId))
	setValue0 := strconv.Itoa(int(relationInfo.UserId))

	ok, err := global.REDIS.SIsMember(ctxx, setKey, setValue).Result()
	global.REDIS.SIsMember(ctxx, setKey0, setValue0).Result()

	if err != nil {
		return ERROR
	}
	if ok == true {
		return ALREADYEXIST
	}

	_, err = global.REDIS.SAdd(ctxx, setKey, setValue).Result()
	global.REDIS.SAdd(ctxx, setKey0, setValue0).Result()
	if err != nil {
		return ERROR
	}
	return SUCCESS
}

//关注后，bitmap将该UserId位置1 方向:ToFollowed->Users
func redisAddRelationByUsers(followInfo model.Follow) int {
	bitmapKey := "relation:" + "followedby_users:" + strconv.Itoa(int(followInfo.FollowedUserId))
	bitmapKey0 := "relation:" + "followby_users:" + strconv.Itoa(int(followInfo.UserId))

	_, err := global.REDIS.SetBit(ctxx, bitmapKey, int64(followInfo.UserId)%4294967296, 1).Result()
	global.REDIS.SetBit(ctxx, bitmapKey0, int64(followInfo.FollowedUserId)%4294967296, 1).Result()
	if err != nil {
		return ERROR
	}
	return SUCCESS
}

//取消赞后，bitmap将该UserId位置0 方向:ToFollowed->Users
func redisDeleteRelationByUsers(followInfo model.Follow) int {
	bitmapKey := "relation:" + "followedby_users:" + strconv.Itoa(int(followInfo.FollowedUserId))
	bitmapKey0 := "relation:" + "followby_users:" + strconv.Itoa(int(followInfo.UserId))

	_, err := global.REDIS.SetBit(ctxx, bitmapKey, int64(followInfo.UserId)%4294967296, 0).Result()
	global.REDIS.SetBit(ctxx, bitmapKey0, int64(followInfo.FollowedUserId)%4294967296, 0).Result()

	if err != nil {
		return ERROR
	}
	return SUCCESS
}

//关注后Redis操作
func (rs *RelationService) RedisAddRelation(followInfo model.Follow) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = redisDeleteUserUnRelation(followInfo); ok == ERROR {
		return false
	}
	if ok = redisAddRelation(followInfo); ok == ERROR {
		return false
	}
	if ok = redisAddRelationByUsers(followInfo); ok == ERROR {
		return false
	}
	return true
}

//取消关注后Redis操作
func (rs *RelationService) RedisDeleteRelation(followInfo model.Follow) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = redisAddUserUnRelations(followInfo); ok == ERROR {
		return false
	}
	if ok = redisUnAddRelation(followInfo); ok == ERROR {
		return false
	}
	if ok = redisDeleteRelationByUsers(followInfo); ok == ERROR {
		return false
	}
	return true
}

func (rs *RelationService) RedisGetFollowList(userId int64) ([]int64, error) {
	var followIds []int64
	zsetKey := "relation:" + "followed:" + strconv.Itoa(int(userId))
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

func (rs *RelationService) RedisGetFollowedList(userId int64) ([]int64, error) {
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

//根据UserId获取用户关注列表
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

//根据UserId获取用户粉丝列表
func (rs *RelationService) GetFollowedList(userId int64) ([]vo.UserVo, error) {

	var followList []model.User
	var followVoList []vo.UserVo

	followIds, err := rs.RedisGetFollowedList(userId)
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
		if ok == BITMAPFOLLOW {
			isDelete = true
		} else {
			isDelete = false
		}

		userVo := vo.UserVo{
			Id:       user.UserId,
			Name:     user.Name,
			IsFollow: isDelete,
		}
		count, _ := RedisGetFollowCount(int64(userVo.Id))
		count1, _ := RedisGetFollowedCount(int64(userVo.Id))
		userVo.FollowCount = uint32(count)
		userVo.FollowerCount = uint32(count1)

		userVos = append(userVos, userVo)
	}
	return userVos
}

func (rs *RelationService) SynchronizeRelationDBAndRedis() {
	log.Println("同步redis到数据库")
	zsetkey, err := global.REDIS.Keys(ctxx, "relation:"+"followed:*").Result()
	log.Println(zsetkey)
	if err != nil {
		return
	}
	var FollowedUserIds []string
	for _, userId := range zsetkey {
		FollowedUserIds, err = global.REDIS.ZRange(ctxx, userId, 0, -1).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		log.Println("like", uid)
		for _, FollowedUserId := range FollowedUserIds {
			log.Println("Followed UserId", FollowedUserId)
			vid := utils.String2Uint64(FollowedUserId)
			rs.AddAction(model.Follow{
				FollowedUserId: vid,
				UserId:         uid,
			})

		}
	}
	var setkey []string
	setkey, err = global.REDIS.Keys(ctxx, "relation:"+"unfollowed:*").Result()
	if err != nil {
		return
	}
	log.Println("unfollow setkey: ", setkey)
	var DeleteFollowedUserIds []string
	for _, userId := range setkey {
		DeleteFollowedUserIds, err = global.REDIS.SMembers(ctxx, userId).Result()
		uid := utils.String2Uint64(utils.SplitString(userId, ":"))
		log.Println("unlike:", uid)
		for _, followeduserId := range DeleteFollowedUserIds {
			log.Println("followed userId:", followeduserId)
			vid := utils.String2Uint64(followeduserId)
			rs.RelationAction(model.Follow{
				UserId:         uid,
				FollowedUserId: vid,
			})
		}

	}
}

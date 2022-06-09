package service

import (
	"context"
	"github.com/DouYin/common/codes"
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/cache"
	"github.com/DouYin/service/repository"
	"log"
)

var (
	followRepository repository.FollowRepository
	ctxx             = context.Background()
)

type RelationService struct {
	relationCache cache.RelationCache
}

func (rs *RelationService) FollowerList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followedUsers, _ := followRepository.GetFollowedOrFollowUserWithUserId(userId, codes.Followed)
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

// RedisAddRelation 关注后Redis操作
func (rs *RelationService) RedisAddRelation(followInfo model.Follow) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = rs.relationCache.RedisDeleteUserUnRelation(followInfo); ok == codes.ERROR {
		return false
	}
	if ok = rs.relationCache.RedisAddRelation(followInfo); ok == codes.ERROR {
		return false
	}
	return true
}

// RedisDeleteRelation 取消关注后Redis操作
func (rs *RelationService) RedisDeleteRelation(followInfo model.Follow) bool {
	var ok int
	//后面加锁，保证原子性
	if ok = rs.relationCache.RedisAddUserUnRelations(followInfo); ok == codes.ERROR {
		return false
	}
	if ok = rs.relationCache.RedisUnAddRelation(followInfo); ok == codes.ERROR {
		return false
	}
	return true
}

// GetFollowList 根据UserId获取用户关注列表
func (rs *RelationService) GetFollowList(userId int64) ([]vo.UserVo, error) {
	var followList []model.User
	var followVoList []vo.UserVo

	followIds, err := rs.relationCache.RedisGetFollowList(userId)
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

	followIds, err := rs.relationCache.RedisGetFollowerList(userId)
	if err != nil {
		return followVoList, err
	}
	for _, id := range followIds {
		followList = append(followList, userRepository.QueryUserDtoInfo(uint64(id)))
	}
	followVoList = rs.FollowList2Vo(userId, followList)
	return followVoList, nil
}

// FollowList 根据UserId获取用户关注列表
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

func (rs *RelationService) FollowList2Vo(userId int64, FollowList []model.User) []vo.UserVo {
	var userVos []vo.UserVo
	for _, user := range FollowList {
		var isDelete bool

		ok, _ := rs.relationCache.RedisIsRelationCreated(int64(user.UserId), userId)
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
		count, _ := rs.relationCache.RedisGetFollowCount(int64(userVo.Id))
		userVo.FollowCount = uint32(count)
		count1, _ := rs.relationCache.RedisGetFollowerCount(int64(userVo.Id))
		userVo.FollowerCount = uint32(count1)

		userVos = append(userVos, userVo)
	}
	return userVos
}

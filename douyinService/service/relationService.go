package service

import (
	"github.com/DouYin/common/codes"
	"log"

	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/request"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/repository"
)

var followRepository repository.FollowRepository

type RelationService struct {
}

func (rs *RelationService) RelationAction(req request.RelationReq, userid uint64) bool {
	where := model.Follow{
		FollowedUserId: req.ToUserId,
		UserId:         userid,
	}
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

func (rs *RelationService) AddAction(req request.RelationReq, userid uint64) bool {
	where := model.Follow{
		FollowedUserId: req.ToUserId,
		UserId:         userid,
	}

	var out model.Follow
	if isOk := followRepository.UpdateFollowUserId(where, &out); !isOk {
		follow := model.Follow{
			UserId:         userid,
			FollowedUserId: req.ToUserId,
			IsDeleted:      false,
		}
		if IsOk := followRepository.AddFollow(follow); !IsOk {
			return false
		}
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
				IsFollow:      user.FollowedA,
			}
			userVos = append(userVos, userVo)
		}
	}
	return userVos
}

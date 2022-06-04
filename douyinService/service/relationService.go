package service

import (
	"github.com/DouYin/common/constant"
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

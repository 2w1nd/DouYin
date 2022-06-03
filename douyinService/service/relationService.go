package service

import (
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
	follow := model.Follow{
		UserId:         userid,
		FollowedUserId: req.ToUserId,
		IsDeleted:      false,
	}
	if isOk := followRepository.AddFollow(follow); !isOk {
		return false
	}
	return true
}

func (rs *RelationService) FollowList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followUsers, _ := followRepository.GetFollowUserWithUserId(userId)
	log.Println(followUsers)
	for _, user := range followUsers {
		log.Println(user.UserId)
	}
	userList = rs.userDto2UserVos(followUsers)
	return userList
}

func (rs *RelationService) FollowerList(userId uint64) []vo.UserVo {
	var userList []vo.UserVo
	followedUsers, _ := followRepository.GetFollowedUserWithUserId(userId)
	log.Println(followedUsers)
	for _, user := range followedUsers {
		log.Println(user.UserId)
	}
	userList = rs.userDto2UserVos(followedUsers)
	return userList
}

func (rs *RelationService) userDto2UserVos(followerUsers []dto.FollowDto) []vo.UserVo {
	var userVos []vo.UserVo
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
	return userVos
}

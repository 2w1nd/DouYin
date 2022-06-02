package service

import (
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/service/repository"
	"github.com/gin-gonic/gin"
	"log"
)

var followRepository repository.FollowRepository

type RelationService struct {
}

func (rs *RelationService) RelationAction(c *gin.Context) {

}

func (rs *RelationService) FollowList(c *gin.Context) {

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

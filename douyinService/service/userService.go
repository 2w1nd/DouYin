package service

import (
	"github.com/DouYin/service/middleware"
	"log"
	"strconv"

	"github.com/DouYin/common/model"
	"github.com/DouYin/service/utils"
	"gorm.io/gorm"
)

type UserService struct {
}

type UserMsg struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint32 `json:"follow_count"`
	FollowerCount uint32 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func (us *UserService) Register(user *model.User) (code int, msg, token string, userId uint64) {
	userbak, err := userRepository.GetUserByUserName(user.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 500, err.Error(), "", 0
	}
	if userbak.Name != "" {
		return 500, "用户已存在", "", 0
	}
	// 用户名不存在，可以注册
	// 密码加盐加密
	saltpassword, salt := utils.Encrypt(user.Password)
	user.Password = saltpassword
	user.Salt = salt
	user, err = userRepository.CreateUser(user)
	if err != nil {
		return 500, err.Error(), "", 0
	}

	// 注册成功，返回用户id和权限jwt token
	token, err = middleware.CreateToken(user.UserId, user.Username)
	if err != nil {
		return 500, err.Error(), "", 0
	}
	return 0, "注册成功", token, user.UserId
}

func (us *UserService) Login(username, password string) (code int, msg, token string, userId uint64) {
	user, err := userRepository.GetUserByUserName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 500, "用户名不存在", "", 0
		}
		return 500, err.Error(), "", 0
	}
	if user.Username == "" {
		return 500, "用户不存在", "", 0
	}

	if user.Password == utils.Analysis(password, user.Salt) {
		// 登录成功，返回用户id和权限jwt token
		token, err = middleware.CreateToken(user.UserId, user.Username)
		if err != nil {
			return 500, err.Error(), "", 0
		}
		return 0, "登陆成功", token, user.UserId
	}
	return 500, "密码错误", "", 0
}

func (us *UserService) UserInfo(id string, myuserId string) (code int, msg string, userMsg UserMsg) {
	userId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("类型转化失败")
		return
	}
	where := model.User{UserId: uint64(userId)}
	user, err := userRepository.GetFirstUser(where)
	if err != nil {
		return 500, err.Error(), UserMsg{}
	}
	myid, err := strconv.Atoi(myuserId)
	if err != nil {
		return 500, err.Error(), UserMsg{}
	}
	_, err = userRepository.GetFollowByUserId(uint64(myid), user.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, "未关注", UserMsg{
				Id:            user.UserId,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false,
			}
		}
		return 500, err.Error(), UserMsg{}
	}
	return 0, "已关注", UserMsg{
		Id:            user.UserId,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      true,
	}
}

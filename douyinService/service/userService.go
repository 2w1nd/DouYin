package service

import (
	"github.com/DouYin/common/codes"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/middleware"
	"golang.org/x/net/context"
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
		return codes.ERROR, err.Error(), "", 0
	}
	if userbak.Name != "" {
		return codes.ERROR, "用户已存在", "", 0
	}
	// 用户名不存在，可以注册
	// 密码加盐加密
	saltpassword, salt := utils.Encrypt(user.Password)
	user.Password = saltpassword
	user.Salt = salt
	user, err = userRepository.CreateUser(user)
	if err != nil {
		return codes.ERROR, err.Error(), "", 0
	}

	// 注册成功，返回用户id和权限jwt token
	token, err = middleware.CreateToken(user.UserId, user.Username)
	if err != nil {
		return codes.ERROR, err.Error(), "", 0
	}
	return codes.SUCCESS, "注册成功", token, user.UserId
}

func (us *UserService) Login(username, password string) (code int, msg, token string, userId uint64) {
	user, err := userRepository.GetUserByUserName(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return codes.ERROR, "用户名不存在", "", 0
		}
		return codes.ERROR, err.Error(), "", 0
	}
	if user.Username == "" {
		return codes.ERROR, "用户不存在", "", 0
	}

	if user.Password == utils.Analysis(password, user.Salt) {
		// 登录成功，返回用户id和权限jwt token
		token, err = middleware.CreateToken(user.UserId, user.Username)
		if err != nil {
			return codes.ERROR, err.Error(), "", 0
		}
		return codes.SUCCESS, "登陆成功", token, user.UserId
	}
	return codes.ERROR, "密码错误", "", 0
}

func (us *UserService) UserInfo(id string, myuserId string) (code int, msg string, userMsg UserMsg) {
	log.Println("查看用户信息")
	userId, err := strconv.Atoi(id)
	if err != nil {
		log.Println("类型转化失败")
		return
	}
	where := model.User{UserId: uint64(userId)}
	user, err := userRepository.GetFirstUser(where)
	if err != nil {
		return codes.ERROR, err.Error(), UserMsg{}
	}
	myid, err := strconv.Atoi(myuserId)
	if err != nil {
		return codes.ERROR, err.Error(), UserMsg{}
	}
	_, err = userRepository.GetFollowByUserId(uint64(myid), user.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			userMsg = UserMsg{
				Id:            user.UserId,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false,
			}
			log.Println("查不到")
			log.Println("user"+myuserId, userMsg.Name)
			global.REDIS.HMSet(context.Background(), "users:user", myuserId, userMsg.Name)
			return codes.SUCCESS, "未关注", userMsg
		}
		return codes.ERROR, err.Error(), UserMsg{}
	}
	userMsg = UserMsg{
		Id:            user.UserId,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      true,
	}
	global.REDIS.HMSet(context.Background(), "users:user", myuserId, userMsg.Name)
	return codes.SUCCESS, "已关注", userMsg
}

func (us *UserService) GetUserName(userId uint64) string {
	where := model.User{UserId: uint64(userId)}
	user, err := userRepository.GetFirstUser(where)
	if err != nil {
		return ""
	}
	return user.Username
}

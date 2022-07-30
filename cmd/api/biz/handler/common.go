// Package handlers common.go: 各 handlers 需要的公共代码
package handler

import (
	user2 "github.com/DouYin/hertz_gen/model/hertz/user"
	user1 "github.com/DouYin/kitex_gen/user"
)

// ========= 请求参数相关 ============

// UserParam 登录/注册时需要获取的参数信息
type UserParam struct {
	UserName string `form:"username" binding:"required,min=2,max=32,alphanumunicode"`
	PassWord string `form:"password" binding:"required,min=5,max=32,alphanumunicode"`
}

// CommonGETParam 大部分需要鉴权的 GET 请求的参数信息
type CommonGETParam struct {
	Uid   int64  `form:"user_id" binding:"required,number"`
	Token string `form:"token" binding:"required,jwt"`
}

// ========= 返回相关 ============

// BaseResponse Gin 返回非预期（错误）结果时使用
type BaseResponse struct {
	Code int64  `json:"status_code"`
	Msg  string `json:"status_msg"`
}

func UserRPC2Hertz(user *user1.User) (usr *user2.User) {
	usr = &user2.User{
		UserId:        user.Id,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      user.IsFollow,
	}
	return
}

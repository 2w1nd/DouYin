package controller

import "github.com/gin-gonic/gin"

// Register
// @Description: 新用户注册时提供用户名，密码，昵称即可，用户名需要保证唯一。创建成功后返回用户id和权限token
// @param: c
func Register(c *gin.Context) {

}

// Login
// @Description: 通过用户名和密码进行登录，登录成功后返回用户id和权限token
// @param: c
func Login(c *gin.Context) {

}

// UserInfo
// @Description: 获取登录用户的id，昵称，如果实现社交部分的功能，还会返回关注数的粉丝数
// @param: c
func UserInfo(c *gin.Context) {

}

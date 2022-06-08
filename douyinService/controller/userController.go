package controller

import (
	"github.com/DouYin/common/codes"
	"log"
	"net/http"
	"time"

	"github.com/DouYin/common/model"
	"github.com/DouYin/service/service"
	"github.com/gin-gonic/gin"
)

var userService service.UserService

// Register
// @Description: 新用户注册时提供用户名，密码，昵称即可，用户名需要保证唯一。创建成功后返回用户id和权限token
// @param: c
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//name := c.Query("name")
	if username == "" || len(username) > 32 || password == "" || len(password) > 32 {
		rlHandler(c, codes.ERROR, "参数错误", "", 0)
		return
	}
	user := &model.User{
		Username: username,
		Name:     username,
		Password: password,
		UserId:   uint64(time.Now().UnixNano()),
	}

	code, msg, token, userId := userService.Register(user)
	rlHandler(c, code, msg, token, userId)
}

// Login
// @Description: 通过用户名和密码进行登录，登录成功后返回用户id和权限token
// @param: c
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || len(username) > 256 || password == "" {
		rlHandler(c, codes.ERROR, "参数错误", "", 0)
		return
	}

	code, msg, token, userId := userService.Login(username, password)
	rlHandler(c, code, msg, token, userId)
}

// UserInfo
// @Description: 获取登录用户的id，昵称，如果实现社交部分的功能，还会返回关注数的粉丝数
// @param: c
func UserInfo(c *gin.Context) {
	log.Println("用户信息接口")
	userId := c.Query("user_id")
	userMsg := service.UserMsg{}
	if userId == "" {
		msgHandler(c, codes.ERROR, "参数错误", &userMsg)
		return
	}
	code, msg, userMsg := userService.UserInfo(userId, c.GetHeader("userId"))
	msgHandler(c, code, msg, &userMsg)
}

func rlHandler(c *gin.Context, statusCode int, msg, token string, userId uint64) {
	code := http.StatusInternalServerError
	if statusCode != codes.ERROR {
		code = http.StatusOK
	}
	c.JSON(code, gin.H{
		"status_code": statusCode,
		"status_msg":  msg,
		"user_id":     userId,
		"token":       token,
	})
}

func msgHandler(c *gin.Context, statusCode int, msg string, userMsg *service.UserMsg) {
	code := http.StatusInternalServerError
	if statusCode != codes.ERROR {
		code = http.StatusOK
	}
	c.JSON(code, gin.H{
		"status_code": statusCode,
		"status_msg":  msg,
		"user":        *userMsg,
	})
}

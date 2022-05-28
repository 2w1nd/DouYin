package controller

import (
	"log"
	"time"

	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userMsg struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint32 `json:"follow_count"`
	FollowerCount uint32 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// Register
// @Description: 新用户注册时提供用户名，密码，昵称即可，用户名需要保证唯一。创建成功后返回用户id和权限token
// @param: c
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	name := c.Query("name")
	if username == "" || len(username) > 32 || password == "" || len(password) > 32 || name == "" || len(name) > 32 {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  "参数错误",
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	user := model.User{}
	err := global.DB.Raw("SELECT name FROM douyin_user WHERE username=?", username).Scan(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  err.Error(),
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	if user.Name != "" {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  "用户已存在",
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	// 用户名不存在，可以注册
	// 密码加盐加密
	saltpassword, salt := utils.Encrypt(password)

	user = model.User{
		Username: username,
		Password: saltpassword,
		Name:     name,
		Salt:     salt,
		UserId:   uint64(time.Now().UnixNano()),
	}
	err = global.DB.Raw("INSERT INTO douyin_user(user_id,name,username,password,salt) VALUES(?,?,?,?,?);",
		user.UserId, user.Name, user.Username, user.Password, user.Salt).Scan(&user).Error
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  err.Error(),
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	// 注册成功，返回用户id和权限jwt token
	token, err := utils.CreateToken(user.UserId, user.Username)
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  err.Error(),
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"status_code":    0,
		"status_message": "注册成功",
		"user_id":        user.UserId,
		"token":          token,
	})
}

// Login
// @Description: 通过用户名和密码进行登录，登录成功后返回用户id和权限token
// @param: c
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	if username == "" || len(username) > 256 || password == "" {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  "参数错误",
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	user := model.User{}
	log.Println(username, password)
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(500, gin.H{
				"status_code": 500,
				"status_msg":  "用户名不存在",
				"user_id":     nil,
				"token":       nil,
			})
			return
		}
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  err.Error(),
			"user_id":     nil,
			"token":       nil,
		})
		return
	}
	if user.Username == "" {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  "用户不存在",
			"user_id":     nil,
			"token":       nil,
		})
		return
	}

	if user.Password == utils.Analysis(password, user.Salt) {
		// 登录成功，返回用户id和权限jwt token
		token, err := utils.CreateToken(user.UserId, user.Username)
		if err != nil {
			c.JSON(500, gin.H{
				"status_code": 500,
				"status_msg":  err.Error(),
				"user_id":     nil,
				"token":       nil,
			})
			return
		}
		c.JSON(200, gin.H{
			"status_code":    0,
			"status_message": "登陆成功",
			"user_id":        user.UserId,
			"token":          token,
		})
		return
	}
	c.JSON(500, gin.H{
		"status_code": 500,
		"status_msg":  "密码错误",
		"user_id":     nil,
		"token":       nil,
	})
}

// UserInfo
// @Description: 获取登录用户的id，昵称，如果实现社交部分的功能，还会返回关注数的粉丝数
// @param: c
func UserInfo(c *gin.Context) {
	username := c.Query("username")
	if username == "" || len(username) > 32 {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  "参数错误",
			"user":        userMsg{},
		})
		return
	}
	user := model.User{}
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  err.Error(),
			"user":        userMsg{},
		})
		return
	}
	follow := model.Follow{}
	err = global.DB.Where("user_id = ? AND followed_user_id = ?", c.GetHeader("user_id"), user.UserId).First(&follow).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(200, gin.H{
				"status_code": 0,
				"status_msg":  "没有关注",
				"user": userMsg{
					Id:            user.UserId,
					Name:          user.Name,
					FollowCount:   user.FollowCount,
					FollowerCount: user.FollowerCount,
					IsFollow:      false,
				},
			})
			return
		}
		c.JSON(500, gin.H{
			"status_code": 500,
			"status_msg":  err.Error(),
			"user":        userMsg{},
		})
		return
	}
	c.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "查询用户信息成功",
		"user": userMsg{
			Id:            user.UserId,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true,
		},
	})
}

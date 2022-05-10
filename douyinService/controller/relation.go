package controller

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RelationAction
// @Description: 登录用户对其他用户进行关注或取消关注
// @param: c
func RelationAction(c *gin.Context) {

}

// FollowList
// @Description: 登录用户关注的所有用户列表
// @param: c
func FollowList(c *gin.Context) {

}

// FollowerList
// @Description: 所有关注登录用户的粉丝列表
// @param: c
func FollowerList(c *gin.Context) {
	DemoUser := vo.UserVo{
		Id:   1,
		Name: "w1nd",
	}
	c.JSON(http.StatusOK, vo.UserListVo{
		Response: response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		UserList: []vo.UserVo{DemoUser},
	})
}

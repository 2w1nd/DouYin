package controller

import (
	"github.com/DouYin/common/entity/request"
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/service"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

var relationService service.RelationService

// RelationAction
// @Description: 登录用户对其他用户进行关注或取消关注
// @param: c
func RelationAction(c *gin.Context) {
	var relationReq request.RelationReq
	_ = c.ShouldBindQuery(&relationReq)
	user := utils.GetUserContext(c)
	where := model.Follow{
		UserId:         user.Id,
		FollowedUserId: relationReq.ToUserId,
	}
	if user.Id == relationReq.ToUserId {
		response.FailWithMessage("不能关注自己噢~", c)
		return
	}
	if relationReq.ActionType == 2 {
		if err := relationService.RedisDeleteRelation(where); !err {
			response.FailWithMessage("取消关注失败", c)
		} else {
			response.OkWithMessage("取消关注成功", c)
		}
	} else if relationReq.ActionType == 1 {
		if err := relationService.RedisAddRelation(where); !err {
			response.FailWithMessage("关注失败", c)
		} else {
			response.OkWithMessage("关注成功", c)
			return
		}
	}
}

// FollowList
// @Description: 登录用户关注的所有用户列表
// @param: c
func FollowList(c *gin.Context) {
	userId := utils.String2Uint64(c.Query("user_id"))
	userList, err := relationService.GetFollowList(int64(userId))
	if err != nil {
		response.FailWithMessage("获取失败", c)
	}
	c.JSON(http.StatusOK, vo.UserListVo{
		Response: response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		UserList: userList,
	})
}

// FollowerList
// @Description: 所有关注登录用户的粉丝列表
// @param: c
func FollowerList(c *gin.Context) {
	userId := utils.String2Uint64(c.Query("user_id"))
	userList, err := relationService.GetFollowerList(int64(userId))
	if err != nil {
		response.FailWithMessage("获取失败", c)
	}
	c.JSON(http.StatusOK, vo.UserListVo{
		Response: response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		UserList: userList,
	})
}

package controller

import (
	"github.com/DouYin/common/entity/dto"
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/service/service"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var publishService service.PublishService

// Publish
// @Description: 登录用户选择视频上传
// @param: c
func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	if err != nil {
		response.FailWithMessage("获取数据失败", c)
		return
	}
	publishService.Publish(utils.GetUserContext(c), data, c.PostForm("title"))
	response.Ok(c)
}

// PublishList
// @Description: 登录用户的视频发布列表，直接列出用户所有投稿过的视频
// @param: c
func PublishList(c *gin.Context) {
	dstUserId, err := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if err != nil {
		response.FailWithMessage("参数user_id有误", c)
		return
	}
	publishList := publishService.PublishList(utils.GetUserContext(c), dstUserId)
	c.JSON(http.StatusOK, dto.VideoListDto{
		Response:  response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		VideoList: publishList,
	})
}

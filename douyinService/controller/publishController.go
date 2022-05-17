package controller

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/service/context"
	"github.com/DouYin/service/service"
	"github.com/gin-gonic/gin"
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
	publishService.Publish(context.UserContext{}, data)
	response.Ok(c)
}

// PublishList
// @Description: 登录用户的视频发布列表，直接列出用户所有投稿过的视频
// @param: c
func PublishList(c *gin.Context) {

}

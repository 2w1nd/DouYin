package controller

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/service/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var feedService service.FeedService

// Feed
// @Description: 无需登录，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
// @param: c
func Feed(c *gin.Context) {

	token := c.Query("token")
	latestTime := c.Query("latest_time")

	video := feedService.Feed(token, latestTime)

	c.JSON(http.StatusOK, vo.VideoListVo{
		Response:  response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		NextTime:  time.Now().Unix(),
		VideoList: video,
	})
}

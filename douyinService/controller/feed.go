package controller

import (
	"github.com/DouYin/common/entity/response"
	model2 "github.com/DouYin/common/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	response.Response
	VideoList []model2.Video `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

// Feed
// @Description: 无需登录，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
// @param: c
func Feed(c *gin.Context) {
	c.JSON(http.StatusOK, FeedResponse{
		Response:  response.Response{StatusCode: 0},
		VideoList: model2.DemoVideos,
		NextTime:  time.Now().Unix(),
	})
}

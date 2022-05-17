package controller

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/service/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var feedService service.FeedService

// Feed
// @Description: 无需登录，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
// @param: c
func Feed(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID := uint64(claims["id"].(float64))
	log.Println(userID)
	token := c.Query("token")
	latestTime := c.Query("latest_time")

	videoList := feedService.Feed(token, latestTime)

	c.JSON(http.StatusOK, vo.VideoListVo{
		Response:  response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		NextTime:  time.Now().Unix(),
		VideoList: videoList,
	})
}

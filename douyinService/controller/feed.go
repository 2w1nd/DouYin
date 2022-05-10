package controller

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Feed
// @Description: 无需登录，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
// @param: c
func Feed(c *gin.Context) {

	DemoUser := model.User{
		Id:            1,
		Name:          "TestUser",
		FollowCount:   0,
		FollowerCount: 0,
	}

	DemoVideos := []vo.VideoVo{
		{
			Id:            1,
			Author:        DemoUser,
			PlayUrl:       "https://www.w3schools.com/html/movie.mp4",
			CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
			FavoriteCount: 0,
			CommentCount:  0,
		},
	}

	c.JSON(http.StatusOK, vo.VideoListVo{
		Response:  response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		NextTime:  1,
		VideoList: DemoVideos,
	})
}

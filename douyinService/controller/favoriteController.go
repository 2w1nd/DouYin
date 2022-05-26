package controller

import (
	"fmt"
	"github.com/DouYin/common/context"
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/service"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"net/http"
)

var favoriteService service.FavoriteService

const (
	LIKE   bool = false
	UNLIKE bool = true
)

type QueryInfo struct {
	UserId     int64  `form:"user_id" binding:"required"`
	Token      string `form:"token" binding:"required"`
	VideoId    int64  `form:"video_id" binding:"required"`
	ActionType int32  `form:"action_type" binding:"required"`
}

// FavoriteAction
// @Description: 登录用户对视频的点赞和取消点赞操作
// @param: c
func FavoriteAction(c *gin.Context) {
	//鉴权
	var userContext context.UserContext
	userContext = utils.GetUserContext(c)
	userID := userContext.Id

	var queryInfo QueryInfo
	err := c.ShouldBindQuery(&queryInfo)
	if err != nil {
		log.Errorf("获取favoriteInfo失败", err)
	}
	fmt.Println("queryInfo: ", queryInfo)

	var favoriteInfo model.Favorite
	var isDeleted bool

	favoriteInfo = model.Favorite{
		UserId:    userID,
		VideoId:   uint64(queryInfo.VideoId),
		IsDeleted: isDeleted,
	}
	//1是点赞(false),2是取消(true)
	if queryInfo.ActionType == 1 {
		favoriteInfo.IsDeleted = LIKE
	} else if queryInfo.ActionType == 2 {
		favoriteInfo.IsDeleted = UNLIKE
	}

	//存入Redis
	ok := favoriteService.RedisAddFavorite(favoriteInfo)

	//如果存入成功返回
	if ok == true {
		c.JSON(-1, response.Response{http.StatusOK, "点赞失败"})
	} else {
		c.JSON(0, response.Response{http.StatusOK, "点赞成功"})
	}
	//c.JSON(0, response.Response{http.StatusOK, "点赞成功"})
}

// FavoriteList
// @Description: 登录用户的所有点赞视频
// @param: c
func FavoriteList(c *gin.Context) {
	//鉴权
	var userContext context.UserContext
	userContext = utils.GetUserContext(c)
	userID := userContext.Id
	log.Infof("userID", userID)

	//在Redis中查询并返回
	videoIds, err := favoriteService.RedisGetFavoriteList(int64(userID))
	if err != nil {
		log.Errorf("RedisGetFavoriteList faild")
	}
	favoriteService.GetFavoriteList(videoIds)

}

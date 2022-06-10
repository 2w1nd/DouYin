package controller

import (
	"fmt"
	"github.com/DouYin/common/entity/request"
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/service"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"strconv"
)

var favoriteService service.FavoriteService

const (
	LIKE   bool = false
	UNLIKE bool = true
)

// FavoriteAction
// @Description: 登录用户对视频的点赞和取消点赞操作
// @param: c
func FavoriteAction(c *gin.Context) {

	var favoriteReq request.FavoriteReq
	err := c.ShouldBindQuery(&favoriteReq)
	if err != nil {
		log.Errorf("获取favoriteReq失败", err)
	}
	fmt.Println("favoriteReq: ", favoriteReq)

	user := utils.GetUserContext(c)

	var favoriteInfo model.Favorite
	var isDeleted bool

	favoriteInfo = model.Favorite{
		UserId:    user.Id,
		VideoId:   uint64(favoriteReq.VideoId),
		IsDeleted: isDeleted,
	}
	//1是点赞(false),2是取消(true)
	if favoriteReq.ActionType == 1 {
		favoriteInfo.IsDeleted = LIKE
	} else if favoriteReq.ActionType == 2 {
		favoriteInfo.IsDeleted = UNLIKE
	}

	var ok bool
	if favoriteInfo.IsDeleted == LIKE {
		//存入Redis
		ok = favoriteService.RedisAddFavorite(favoriteInfo)
		if ok == true {
			c.JSON(0, response.Response{StatusCode: response.SUCCESS, StatusMsg: "点赞成功"})
		} else {
			c.JSON(-1, response.Response{StatusCode: response.ERROR, StatusMsg: "点赞失败"})
		}
	} else {
		//从Redis中删除
		ok = favoriteService.RedisDeleteFavorite(favoriteInfo)
		if ok == true {
			c.JSON(0, response.Response{StatusCode: response.SUCCESS, StatusMsg: "取消赞成功"})
		} else {
			c.JSON(-1, response.Response{StatusCode: response.ERROR, StatusMsg: "取消赞失败"})
		}
	}

}

// FavoriteList
// @Description: 登录用户的所有点赞视频
// @param: c
func FavoriteList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	fmt.Println("userId:", userId)

	if err != nil {
		c.JSON(response.ERROR, vo.FavoriteListVo{
			Response: response.Response{StatusCode: -1, StatusMsg: "操作失败"},
		})
	}
	//在Redis中查询并返回
	videos, err := favoriteService.GetFavoriteList(userId)
	if err != nil {
		c.JSON(response.ERROR, vo.FavoriteListVo{
			Response: response.Response{StatusCode: -1, StatusMsg: "操作失败"},
		})
	}

	c.JSON(response.SUCCESS, vo.FavoriteListVo{
		Response:  response.Response{StatusMsg: "操作成功"},
		VideoList: videos,
	})

}

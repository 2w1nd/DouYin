package controller

import (
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/service"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"net/http"
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
	//鉴权
	claims := jwt.ExtractClaims(c)
	userID := uint64(claims["id"].(float64))
	log.Infof("userID", userID)

	var err error
	var favoriteInfo model.Favorite

	err = c.ShouldBindQuery(&favoriteInfo)
	if err != nil {
		log.Errorf("获取favoriteInfo失败", err)
	}
	//log.Infof("favoriteInfo", favoriteInfo)
	var isDeleted int64
	isDeleted, err = strconv.ParseInt(c.Query("action_type"), 10, 64)
	if err != nil {
		log.Errorf("获取action_type失败", err)
	}

	//1是点赞(false),2是取消(true)
	if isDeleted == 1 {
		favoriteInfo.IsDeleted = LIKE
	} else if isDeleted == 2 {
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

}

// FavoriteList
// @Description: 登录用户的所有点赞视频
// @param: c
func FavoriteList(c *gin.Context) {
	//鉴权
	claims := jwt.ExtractClaims(c)
	userID := uint64(claims["id"].(float64))
	log.Infof("userID", userID)

	//在Redis中查询并返回
	videoIds, err := favoriteService.RedisGetFavoriteList(int64(userID))
	if err != nil {
		log.Errorf("RedisGetFavoriteList faild")
	}
	favoriteService.GetFavoriteList(videoIds)

}

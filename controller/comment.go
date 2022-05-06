package controller

import (
	"github.com/DouYin/model"
	_ "github.com/DouYin/model"
	response "github.com/DouYin/model/reponse"
	"github.com/DouYin/service"
	"github.com/gin-gonic/gin"
	"log"
)

type CommentController struct {
}

type CommentListResponse struct {
	response.Response
	CommentList []response.Response `json:"comment_list,omitempty"`
}

// AddCommentDemo
// @Description: 添加评论到数据库（栗子）
// @receiver: e
// @param: c
func AddCommentDemo(c *gin.Context) {
	var comment model.Comment
	_ = c.ShouldBindJSON(&comment)
	log.Println(comment.Content)
	if err := service.AddCommentDemo(comment); err != nil {
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// CommentAction
// @Description: 登录用户对视频进行评论
// @param: c
func CommentAction(c *gin.Context) {
	var userLoginInfo map[string]model.User
	_ = c.ShouldBindJSON(&userLoginInfo)

	token := c.Query("token")

	if _, exist := userLoginInfo[token]; exist {
		response.OkWithMessage("成功", c)
	} else {
		response.FailWithMessage("错误", c)
	}
}

// CommentList
// @Description: 查看视频的所有评论，按发布时间倒序
// @param: c
func CommentList(c *gin.Context) {
	response.OkWithMessage("成功", c)
}

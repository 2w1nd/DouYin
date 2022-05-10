package controller

import (
	"github.com/DouYin/common/entity/response"
	model2 "github.com/DouYin/common/model"
	"github.com/DouYin/service/service"
	"github.com/gin-gonic/gin"
	"log"
)

var commentService service.CommentService

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
	var comment model2.Comment
	_ = c.ShouldBindJSON(&comment)
	log.Println(comment.Content)
	if err := commentService.AddCommentDemo(comment); err != nil {
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// CommentAction
// @Description: 登录用户对视频进行评论
// @param: c
func CommentAction(c *gin.Context) {
	var userLoginInfo map[string]model2.User
	_ = c.ShouldBindJSON(&userLoginInfo)

	token := c.Query("token")

	// TODO 调用service层函数进行处理

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

	// TODO 调用service层函数进行处理

	response.OkWithMessage("成功", c)
}

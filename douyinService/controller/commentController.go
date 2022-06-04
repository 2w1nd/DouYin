package controller

import (
	"github.com/DouYin/common/entity/request"
	"github.com/DouYin/common/entity/response"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/service"
	"github.com/DouYin/service/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

var commentService service.CommentService

// AddCommentDemo
// @Description: 添加评论到数据库（栗子）
// @receiver: e
// @param: c
func AddCommentDemo(c *gin.Context) {
	var comment model.Comment
	_ = c.ShouldBindQuery(&comment)
	comment.CommentId = uint64(global.ID.Generate())
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
	var commentReq request.CommentReq
	_ = c.ShouldBindQuery(&commentReq)
	if commentReq.ActionType == 1 {
		validate := validator.New()
		err := validate.Struct(commentReq)
		//添加评论
		if err != nil {
			response.FailWithMessage("评论至少为2个文字", c)
			// from here you can create your own error messages in whatever language you wish
			return
		} else {
			if CommentVos, err := commentService.AddComment(commentReq, utils.GetUserContext(c)); !err {
				response.FailWithMessage("评论失败", c)
			} else {
				c.JSON(http.StatusOK, vo.CommentRet{
					Response: response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
					Comment:  CommentVos[0],
				})
			}
		}
	} else if commentReq.ActionType == 2 {
		// 删除评论
		if isOk := commentService.DeleteComment(commentReq); !isOk {
			response.FailWithMessage("删除失败", c)
		} else {
			response.OkWithMessage("删除成功", c)
		}
	} else {
		response.FailWithMessage("操作失败", c)
	}
}

// CommentList
// @Description: 查看视频的所有评论，按发布时间倒序
// @param: c
func CommentList(c *gin.Context) {
	videoId, _ := strconv.ParseUint(c.Query("video_id"), 10, 64)
	commentVos := commentService.GetCommentList(videoId)
	c.JSON(http.StatusOK, vo.CommentListVo{
		Response:    response.Response{StatusCode: response.SUCCESS, StatusMsg: "操作成功"},
		CommentList: commentVos,
	})
}

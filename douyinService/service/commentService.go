package service

import (
	"github.com/DouYin/common/context"
	"github.com/DouYin/common/entity/request"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
)

type CommentService struct {
}

var commentRepository repository.CommentRepository

// AddCommentDemo
// @Description: 测试栗子
// @receiver: e
// @param: c
// @return: err
func (cs *CommentService) AddCommentDemo(c model.Comment) (err error) {
	err = global.DB.Create(&c).Error
	return err
}

func (cs *CommentService) AddComment(req request.CommentReq, context context.UserContext) ([]vo.CommentVo, bool) {
	var (
		commentVos  []vo.CommentVo
		commentList []model.Comment
	)
	userId := context.Id
	commentId := uint64(global.ID.Generate())
	comment := model.Comment{
		CommentId: commentId,
		UserId:    userId,
		VideoId:   req.VideoId,
		Content:   req.CommentTest,
	}
	isOk := commentRepository.AddComment(comment)
	if !isOk {
		return commentVos, false
	}
	commentRet, _ := commentRepository.QueryCommentWithUserInfo(commentId)
	commentList = append(commentList, commentRet)
	return cs.commentList2Vo(commentList, req.VideoId), true
}

func (cs *CommentService) DeleteComment(req request.CommentReq) bool {
	where := model.Comment{CommentId: req.CommentId}
	if isOk := commentRepository.DeleteCommentById(where); !isOk {
		return false
	}
	return true
}

func (cs *CommentService) GetCommentList(userContext context.UserContext, videoId uint64) []vo.CommentVo {
	var commentVos []vo.CommentVo
	commentList, _ := commentRepository.CommentListByVideoId(videoId)
	commentVos = cs.commentList2Vo(commentList, videoId)
	return commentVos
}

func (cs *CommentService) commentList2Vo(CommentList []model.Comment, videoID uint64) []vo.CommentVo {
	var commentVos []vo.CommentVo
	video := videoRepository.GetVideoByVideoId(videoID)
	for _, comment := range CommentList {
		var isFollow []model.Follow
		query := global.DB.Debug().Model(model.Follow{}).Where(" user_id=? and followed_user_id=?", comment.UserId, video.AuthorId)
		query.Find(&isFollow)
		var isFollowed bool
		if len(isFollow) != 0 {
			isFollowed = true
		} else {
			isFollowed = false
		}
		commentVo := vo.CommentVo{
			CommentId: comment.CommentId,
			User: vo.CommentUserVo{
				UserID:        comment.CommentUser.UserId,
				Name:          comment.CommentUser.Username,
				FollowCount:   comment.CommentUser.FollowCount,
				FollowerCount: comment.CommentUser.FollowerCount,
				IsFollow:      isFollowed,
			},
			Content:    comment.Content,
			CreateDate: comment.GmtCreated.Format("01-02"),
		}
		commentVos = append(commentVos, commentVo)
	}
	return commentVos
}

package repository

import (
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
)

type CommentRepository struct {
	Base BaseRepository
}

func (cr *CommentRepository) AddComment(comment model.Comment) bool {
	if err := cr.Base.Create(&comment); err != nil {
		return false
	}
	var preComment model.Video
	query := global.DB.Model(model.Video{}).
		Where("video_id =?", comment.VideoId)
	query.Find(&preComment)
	global.DB.Model(model.Video{}).
		Where("video_id = ?", comment.VideoId).
		Update("comment_count", preComment.CommentCount+1)
	return true
}

func (cr *CommentRepository) DeleteCommentById(where interface{}, id uint64) bool {
	var comment model.Comment
	if err := cr.Base.DeleteSoftByID(where, &comment); err != nil {
		return false
	}
	var preComment model.Video
	query := global.DB.Model(model.Video{}).
		Where("video_id=?", id)
	query.Find(&preComment)
	global.DB.Model(model.Video{}).Debug().
		Where("video_id=?", id).
		Update("comment_count", preComment.CommentCount-1)
	return true
}

func (cr *CommentRepository) QueryCommentWithUserInfo(commentId uint64) (model.Comment, bool) {
	var comment model.Comment
	query := global.DB.
		Model(model.Comment{}).
		Where("comment_id = ?", commentId).
		Order("gmt_created desc").
		Preload("CommentUser")
	query.Find(&comment)
	return comment, true
}

func (cr *CommentRepository) CommentListByVideoId(videoId uint64) ([]model.Comment, int64) {
	var commentList []model.Comment

	query := global.DB.Debug().
		Model(model.Comment{}).
		Where("video_id = ? and is_deleted != 1", videoId).
		Order("gmt_created desc").
		Preload("CommentUser")
	rows := query.Find(&commentList).RowsAffected
	return commentList, rows
}

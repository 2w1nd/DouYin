package service

import (
	"encoding/json"
	"fmt"
	"github.com/DouYin/common/entity/request"
	"github.com/DouYin/common/entity/vo"
	"github.com/DouYin/common/model"
	"github.com/DouYin/service/global"
	"github.com/DouYin/service/repository"
	"golang.org/x/net/context"
	"log"
	"strconv"
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

func (cs *CommentService) AddComment(req request.CommentReq) ([]vo.CommentVo, bool) {
	var (
		commentVos  []vo.CommentVo
		commentList []model.Comment
	)
	userId := 0
	commentId := uint64(global.ID.Generate())
	comment := model.Comment{
		CommentId: commentId,
		UserId:    uint64(userId),
		VideoId:   req.VideoId,
		Content:   req.CommentTest,
	}
	isOk := commentRepository.AddComment(comment)
	if !isOk {
		return commentVos, false
	}
	commentRet, _ := commentRepository.QueryCommentWithUserInfo(commentId)
	commentList = append(commentList, commentRet)
	global.REDIS.Del(context.Background(), strconv.FormatUint(req.VideoId, 10))
	return cs.commentList2Vo(commentList, req.VideoId), true
}

func (cs *CommentService) DeleteComment(req request.CommentReq) bool {
	where := model.Comment{CommentId: req.CommentId}
	if isOk := commentRepository.DeleteCommentById(where); !isOk {
		return false
	}
	// 放入缓存
	global.REDIS.Del(context.Background(), strconv.FormatUint(req.VideoId, 10))
	return true
}

func (cs *CommentService) GetCommentList(videoId uint64) []vo.CommentVo {
	var commentVos []vo.CommentVo
	CommentString := "视频的ID为"
	fmt.Println(CommentString + strconv.FormatUint(videoId, 10))
	data1, _ := global.REDIS.Get(context.Background(), CommentString+strconv.FormatUint(videoId, 10)).Result()
	fmt.Println(data1)
	if data1 != "" {
		log.Println("从缓存中查询")
		err := json.Unmarshal([]byte(data1), &commentVos)
		if err != nil {
			return nil
		}

	} else {
		// 从数据库中查询
		log.Println("从数据库中查询")
		commentList, _ := commentRepository.CommentListByVideoId(videoId)
		// 放入缓存
		commentVos = cs.commentList2Vo(commentList, videoId)
		data, _ := json.Marshal(commentVos)
		global.REDIS.Set(context.Background(), CommentString+strconv.FormatUint(videoId, 10), data, 0)
	}
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

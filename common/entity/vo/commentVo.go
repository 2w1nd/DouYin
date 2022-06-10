package vo

import (
	"github.com/DouYin/common/entity/response"
)

type CommentListVo struct {
	response.Response             //状态码，0-成功，其他值-失败
	CommentList       []CommentVo `json:"comment_list,omitempty"` //评论列表
}

type CommentRet struct {
	response.Response             //状态码，0-成功，其他值-失败
	Comment CommentVo `json:"comment"`
}

type CommentVo struct {
	CommentId     uint64        `json:"id,omitempty"` //评论id
	User       CommentUserVo `json:"user"`
	Content    string        `json:"content"`     //评论内容
	CreateDate string        `json:"create_date"` //评论发布日期，格式 mm-dd
}

type CommentUserVo struct {
	UserID        uint64 `json:"id,omitempty"`             //用户id
	Name          string `json:"name,omitempty"`           //用户名称
	FollowCount   uint32 `json:"follow_count,omitempty"`   //关注总数
	FollowerCount uint32 `json:"follower_count,omitempty"` //粉丝总数
	IsFollow      bool   `json:"is_follow"`                //true-已关注，false-未关注
}

package vo

import "github.com/DouYin/common/entity/response"

type UserListVo struct {
	response.Response
	UserList []UserVo `json:"user_list"`
}

type UserVo struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

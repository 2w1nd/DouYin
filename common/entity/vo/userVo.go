package vo

import "github.com/DouYin/common/entity/response"

type UserListVo struct {
	response.Response
	UserList []UserVo `json:"user_list"`
}

type UserInfoVo struct {
	response.Response
	User UserVo `json:"user"`
}

type UserVo struct {
	Id            uint64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint32  `json:"follow_count,omitempty"`
	FollowerCount uint32  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type LoginVo struct {
	response.Response
	UserId 		  uint64 `json:"user_id"`
	Token		  string `json:"token"`
}

package dto

type UserDto struct {
	Id            uint64 `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint32 `json:"follow_count,omitempty"`
	FollowerCount uint32 `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

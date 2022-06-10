package dto

type UserDto struct {
	Id            uint64 `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   uint32 `json:"follow_count"`
	FollowerCount uint32 `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

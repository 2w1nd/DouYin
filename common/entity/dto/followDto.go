package dto

type FollowDto struct {
	UserId uint64
	Name  string
	FollowedUserId uint64
	FollowedA bool
	FollowedB bool
	FollowCount uint32
	FollowerCount uint32
}

package request

type RelationReq struct {
	Token     string `json:"token" form:"token"`
	ToUserId  uint64 `json:"to_user_id" form:"to_user_id"` // 被关注者id
	IsDeleted bool   `json:"is_deleted" form:"is_deleted"` // 关注状态，0关注，1取消
}

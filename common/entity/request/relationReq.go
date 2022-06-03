package request

type RelationReq struct {
	Token     string `json:"token" form:"token"`
	ToUserId  uint64 `json:"to_user_id" form:"to_user_id"` // 被关注者id
	ActionType int   `json:"action_type" form:"action_type"` // 关注状态，0关注，1取消
}

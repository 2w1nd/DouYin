package request

type CommentReq struct {
	UserId      uint64 `json:"user_id" form:"user_id"`
	Token       string `json:"token" form:"token"`
	VideoId     uint64 `json:"video_id" form:"video_id"`
	ActionType  int    `json:"action_type" form:"action_type"` // 1-发布评论，2-删除评论
	CommentTest string `json:"comment_text" form:"comment_text" validate:"min=2,max=20"`
	CommentId   uint64 `json:"comment_id" form:"comment_id"`
}

package request

type FavoriteReq struct {
	Token      string `form:"token" binding:"required"`
	VideoId    int64  `form:"video_id" binding:"required"`
	ActionType int32  `form:"action_type" binding:"required"`
}

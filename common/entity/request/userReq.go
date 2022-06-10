package request


type UserReq struct {
	Username string `json:"username" validate:"required,min=3,max=35"`
	Password string `json:"password" validate:"required,min=6,max=35"`
}

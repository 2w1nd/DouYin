package auth

import (
	"context"
	"github.com/DouYin/cmd/api/rpc"
	"github.com/DouYin/hertz_gen/model/hertz/user"
	user1 "github.com/DouYin/kitex_gen/user"
	"github.com/DouYin/pkg/constants"
	"github.com/bytedance/gopkg/util/logger"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"
	"time"
)

type loginResponse struct {
	Code  int64  `json:"status_code"`
	Msg   string `json:"status_msg"`
	Uid   int64  `json:"user_id"`
	Token string `json:"token"`
}

var Config = &jwt.HertzJWTMiddleware{
	Realm:            "DouYin",                       // 标识
	SigningAlgorithm: "HS256",                        // 加密算法
	Key:              []byte(constants.JWTSecretKey), // 密钥
	Timeout:          time.Hour * 24,                 // token的过期时间
	MaxRefresh:       time.Hour,                      // 刷新最大延时
	IdentityKey:      constants.IdentityKey,
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*user.User); ok {
			return jwt.MapClaims{
				constants.IdentityKey: v.UserId,
			}
		}
		return jwt.MapClaims{}
	},
	IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
		logger.Info("IdentityHandler")
		claims := jwt.ExtractClaims(ctx, c)
		return &user.User{
			UserId: int64(claims[constants.IdentityKey].(float64)),
		}
	},
	Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
		var loginVals user.LoginReq
		if err := c.BindAndValidate(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		if uid, err := rpc.CheckUser(ctx, &user1.CheckUserRequest{Username: loginVals.Username, Password: loginVals.Password}); err == nil {
			c.Set("key", uid)
			return &user.User{
				UserId: uid,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	},
	LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, time time.Time) {
		uId, _ := c.Get("key")
		c.JSON(200, loginResponse{
			Code:  200,
			Msg:   "登录成功",
			Uid:   uId.(int64),
			Token: token,
		})
	},
	Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
		if _, ok := data.(*user.User); ok {
			//logger.Info(usr)
			return true
		}
		return false
	},
	Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
		c.JSON(code, map[string]interface{}{
			"code":    code,
			"message": message,
		})
	},
	CookieName:    "token",
	TokenLookup:   "header: Authorization, query: token, cookie: jwt",
	TokenHeadName: "Bearer",
	TimeFunc:      time.Now,
}

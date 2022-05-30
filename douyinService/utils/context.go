package utils

import (
	"net/http"
	"strconv"
	"time"

	"github.com/DouYin/common/context"
	jwt "github.com/appleboy/gin-jwt/v2"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtkey = []byte("JjUhqZteNUhtDQfvXH9uCHhdKDmUDyAm")

type MyClaims struct {
	UserID   uint64 `json:"userID"`
	UserName string `json:"userName"`
	jwtgo.StandardClaims
}

func GetUserContext(c *gin.Context) context.UserContext {
	claims := jwt.ExtractClaims(c)
	if claims["id"] == nil {
		return context.UserContext{Id: 0}
	}
	return context.UserContext{
		Id: uint64(claims["id"].(float64)),
	}
}

// CreateToken 生成token
func CreateToken(userId uint64, username string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour) //过期时间
	nowTime := time.Now()                        //当前时间
	claims := MyClaims{
		UserID:   userId,
		UserName: username,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "Douyin",          //颁发者签名
			Subject:   "userToken",       //签名主题
		},
	}
	tokenStruct := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return tokenStruct.SignedString(jwtkey)
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, bool) {
	tokenObj, _ := jwtgo.ParseWithClaims(token, &MyClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return jwtkey, nil
	})
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

// JwtMiddleware jwt中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从query中获取token
		tokenStr := c.Query("token")
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "用户不存在"})
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := CheckToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "token不正确"})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, gin.H{"status_code": 0, "status_msg": "token过期"})
			c.Abort() //阻止执行
			return
		}
		c.Request.Header.Set("userId", strconv.FormatUint(tokenStruck.UserID, 10))
		c.Request.Header.Set("userName", tokenStruck.UserName)
		c.Next()
	}
}

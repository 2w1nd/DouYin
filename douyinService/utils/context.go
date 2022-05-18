package utils

import (
	"github.com/DouYin/common/context"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetUserContext(c *gin.Context) context.UserContext {
	claims := jwt.ExtractClaims(c)
	if claims["id"] == nil {
		return context.UserContext{Id: 0}
	}
	return context.UserContext{
		Id: uint64(claims["id"].(float64)),
	}
}

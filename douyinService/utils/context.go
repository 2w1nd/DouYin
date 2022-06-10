package utils

import (
	"github.com/DouYin/common/context"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetUserContext(c *gin.Context) context.UserContext {
	userId, userName := c.GetHeader("userId"), c.GetHeader("userName")
	if userId == "" || userName == "" {
		return context.UserContext{Id: 0, Name: ""}
	}
	id, _ := strconv.Atoi(userId)
	return context.UserContext{Id: uint64(id), Name: userName}
}

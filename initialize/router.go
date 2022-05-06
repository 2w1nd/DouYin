package initialize

import (
	_router "github.com/DouYin/router"
	"github.com/gin-gonic/gin"
)

// Routers
// @Description: 初始化总路由
// @return: *gin.Engine
func Routers() *gin.Engine {
	var Router = gin.Default()

	var router _router.Router
	PrivateGroup := Router.Group("")
	PrivateGroup.Use()
	{
		router.InitRouter(PrivateGroup)
	}
	return Router
}

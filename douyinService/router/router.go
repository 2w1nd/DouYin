package router

import (
	"github.com/DouYin/service/controller"
	"github.com/DouYin/service/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
}

func (rt *Router) InitRouter(r *gin.RouterGroup) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	//// basic apis
	apiRouter.GET("/feed/", middleware.NewJWTFuncWithAction(), controller.Feed)
	apiRouter.GET("/user/", middleware.JwtMiddleware(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	apiRouter.POST("/publish/action/", middleware.JwtMiddleware(), controller.Publish)
	apiRouter.GET("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.JwtMiddleware(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", controller.FavoriteList)
	apiRouter.POST("/comment/demo/add/", controller.AddCommentDemo)
	apiRouter.POST("/comment/action/", middleware.JwtMiddleware(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.JwtMiddleware(), controller.CommentList)

	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.JwtMiddleware(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", controller.FollowList)
	apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}

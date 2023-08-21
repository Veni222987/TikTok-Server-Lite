package main

import (
	"DoushengABCD/controller"
	"DoushengABCD/service"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Static("/static", "./public")

	router := r.Group("/douyin")

	// basic apis
	router.GET("/feed/", controller.Feed)
	router.GET("/user/", controller.GetUserInfo)
	router.POST("/user/register/", controller.Register)
	router.POST("/user/login/", controller.Login)
	router.POST("/publish/action/", service.FormAuthMiddleWare(), controller.UploadVideo)
	router.GET("/publish/list/", service.QueryAuthMiddleWare(), controller.PublishList)

	// extra apis - I
	router.POST("/favorite/action/", service.QueryAuthMiddleWare(), controller.Like)
	router.GET("/favorite/list/", service.QueryAuthMiddleWare(), controller.GetFavoriteList)
	router.POST("/comment/action/", service.QueryAuthMiddleWare(), controller.Comment)
	router.GET("/comment/list/", service.QueryAuthMiddleWare(), controller.GetCommentList)

	// extra apis - II
	router.POST("/relation/action/", service.QueryAuthMiddleWare(), controller.RelationAction)
	router.GET("/relation/follow/list/", service.QueryAuthMiddleWare(), controller.FollowList)
	router.GET("/relation/follower/list/", service.QueryAuthMiddleWare(), controller.FollowerList)
	router.GET("/relation/friend/list/", service.QueryAuthMiddleWare(), controller.FriendList)
	router.GET("/message/chat/")
	router.POST("/message/action/")

}

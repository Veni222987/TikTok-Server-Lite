package main

import (
	"DoushengABCD/controller"
	"DoushengABCD/middleware"
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
	router.POST("/publish/action/", middleware.FormAuthMiddleWare(), controller.UploadVideo)
	router.GET("/publish/list/", middleware.QueryAuthMiddleWare(), controller.PublishList)

	// extra apis - I
	router.POST("/favorite/action/", middleware.QueryAuthMiddleWare(), controller.Like)
	router.GET("/favorite/list/", middleware.QueryAuthMiddleWare(), controller.GetFavoriteList)
	router.POST("/comment/action/", middleware.QueryAuthMiddleWare(), controller.Comment)
	router.GET("/comment/list/", middleware.QueryAuthMiddleWare(), controller.GetCommentList)

	// extra apis - II
	router.POST("/relation/action/", middleware.QueryAuthMiddleWare(), controller.RelationAction)
	router.GET("/relation/follow/list/", middleware.QueryAuthMiddleWare(), controller.FollowList)
	router.GET("/relation/follower/list/", middleware.QueryAuthMiddleWare(), controller.FollowerList)
	router.GET("/relation/friend/list/", middleware.QueryAuthMiddleWare(), controller.FriendList)
	router.GET("/message/chat/", middleware.QueryAuthMiddleWare(), controller.GetChatHistory)
	router.POST("/message/action/", middleware.QueryAuthMiddleWare(), controller.SendMessage)

}

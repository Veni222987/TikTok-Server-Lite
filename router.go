package main

import (
	"DoushengABCD/controller"
	"DoushengABCD/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Static("/static", "./public")

	router := r.Group("/douyin")

	// 基础功能
	router.GET("/feed/", controller.Feed)
	router.GET("/user/", controller.GetUserInfo)
	router.POST("/user/register/", controller.Register)
	router.POST("/user/login/", controller.Login)
	router.POST("/publish/action/", middleware.FormAuthMiddleware(), controller.UploadVideo)
	router.GET("/publish/list/", middleware.QueryAuthMiddleware(), controller.PublishList)

	// 互动功能
	router.POST("/favorite/action/", middleware.QueryAuthMiddleware(), controller.Like)
	router.GET("/favorite/list/", middleware.QueryAuthMiddleware(), controller.GetFavoriteList)
	router.POST("/comment/action/", middleware.QueryAuthMiddleware(), controller.Comment)
	router.GET("/comment/list/", middleware.QueryAuthMiddleware(), controller.GetCommentList)

	// 社交功能
	router.POST("/relation/action/", middleware.QueryAuthMiddleware(), controller.RelationAction)
	router.GET("/relation/follow/list/", middleware.QueryAuthMiddleware(), controller.FollowList)
	router.GET("/relation/follower/list/", middleware.QueryAuthMiddleware(), controller.FollowerList)
	router.GET("/relation/friend/list/", middleware.QueryAuthMiddleware(), controller.FriendList)
	router.GET("/message/chat/", middleware.QueryAuthMiddleware(), controller.GetChatHistory)
	router.POST("/message/action/", middleware.QueryAuthMiddleware(), controller.SendMessage)

}

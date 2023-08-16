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
	router.POST("/publish/action/", service.AuthMiddleWare(), controller.UploadVideo)
	router.GET("/publish/list/", service.AuthMiddleWare(), controller.PublishList)

	// extra apis - I
	router.POST("/favorite/action/", service.AuthMiddleWare(), controller.Like)
	router.GET("/favorite/list/", service.AuthMiddleWare(), controller.GetFavoriteList)
	router.POST("/comment/action/", service.AuthMiddleWare(), controller.Comment)
	router.GET("/comment/list/", service.AuthMiddleWare(), controller.GetCommentList)

	// extra apis - II
	router.POST("/relation/action/", service.AuthMiddleWare(), controller.RelationAction)
	router.GET("/relation/follow/list/", service.AuthMiddleWare(), controller.FollowList)
	router.GET("/relation/follower/list/", service.AuthMiddleWare(), controller.FollowerList)
	router.GET("/relation/friend/list/", service.AuthMiddleWare(), controller.FriendList)
	router.GET("/message/chat/")
	router.POST("/message/action/")

}

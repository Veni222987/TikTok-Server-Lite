package main

import (
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	r.Static("/static", "./public")
	router := r.Group("/dousheng")
	router.GET("/feed/")

}

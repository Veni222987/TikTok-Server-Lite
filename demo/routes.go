package main

import (
	"ginEssential/controller"
	"github.com/gin-gonic/gin"
)

func collectRoute(r *gin.Engine) *gin.Engine {
	r.POST("/api/auth/resgister", controller.Register)
	r.POST("/api/auth/login", controller.Login)

	return r
}

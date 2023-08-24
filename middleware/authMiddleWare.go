package middleware

import (
	"DoushengABCD/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

func QueryAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if service.IsTokenExist(token) {
			//fmt.Println("鉴权成功，token有效\n")
			service.RedisClient.Set(token, service.RedisClient.Get(token).Result, 86400000000000)
			ctx.Next()
		} else {
			fmt.Println("无效的token")
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "无效的Token",
			})
			return
		}
	}
}

func FormAuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.PostForm("token")

		if service.IsTokenExist(token) {
			//fmt.Println("鉴权成功，token有效\n")
			service.RedisClient.Set(token, service.RedisClient.Get(token).Result, 86400000000000)
			ctx.Next()
		} else {
			fmt.Println("无效的token")
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "无效的Token",
			})
			return
		}
	}
}

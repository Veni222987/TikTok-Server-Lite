package middleware

import (
	"DoushengABCD/service"
	"fmt"
	"github.com/gin-gonic/gin"
)

// QueryAuthMiddleware 检测Paras里面的token有效性
func QueryAuthMiddleware() gin.HandlerFunc {
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

// FormAuthMiddleware 检测表单token的有效性
func FormAuthMiddleware() gin.HandlerFunc {
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

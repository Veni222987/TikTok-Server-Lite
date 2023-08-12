package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if IsTokenExist(token) {
			fmt.Println("鉴权成功，token有效\n")
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "无效的Token",
			})
			return
		}
	}
}

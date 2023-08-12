package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func Like(ctx *gin.Context) {
	token := ctx.Query("token")
	video_id := ctx.Query("video_id")
	action_type := ctx.Query("action_type")
	vid64, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		panic(err)
	}
	uid64 := service.GetIdByToken(token)
	liekLog := model.Like{UserId: uid64, VideoId: vid64}
	if action_type == "1" {
		//点赞
		//更新like
		model.Db.Create(&liekLog)
		//更新video表和user表
		model.Db.Where("author_id=?", uid64).Update("favorite_count", gorm.Expr("favorite_count+1"))
		model.Db.Where("id=?", uid64).Update("favorite_count", gorm.Expr("favorite_count+1"))
		ctx.JSON(200, gin.H{
			"status_code": 200,
			"status_msg":  "点赞成功",
		})
	} else if action_type == "2" {
		//取消点赞
		//更新like
		model.Db.Delete(&liekLog)
		//更新video表和user表
		model.Db.Where("author_id=?", uid64).Update("favorite_count", gorm.Expr("favorite_count-1"))
		model.Db.Where("id=?", uid64).Update("favorite_count", gorm.Expr("favorite_count-1"))
		ctx.JSON(200, gin.H{
			"status_code": 200,
			"status_msg":  "取消点赞成功",
		})
	} else {
		//非法请求
		ctx.JSON(400, gin.H{
			"status_code": 400,
			"status_msg":  "非法请求",
		})
	}
}

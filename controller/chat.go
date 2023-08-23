package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func SendMessage(ctx *gin.Context) {
	actionType := ctx.Query("action_type")
	if len(ctx.Query("content")) == 0 {
		ctx.JSON(400, gin.H{
			"status_code": 1,
			"status_msg":  "content不能为空",
		})
		return
	}

	if actionType == "1" {
		uid := service.GetIdByToken(ctx.Query("token"))
		toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
		if err != nil {
			panic(err)
		}
		content := ctx.Query("content")
		//插入数据库
		messageLog := model.Message{OriginUserId: uid, DestinationUserId: toUserId, Content: content, CreateDate: time.Now().UnixNano() / 1e6}
		if res := model.Db.Create(&messageLog); res.Error != nil {
			panic(res.Error)
		}

		ctx.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "发送成功",
		})
	} else {
		ctx.JSON(400, gin.H{
			"status_code": 1,
			"status_msg":  "actionType错误",
		})
	}
}

// 获取聊天记录
func GetChatHistory(ctx *gin.Context) {
	//println(ctx.Query("token"), "最新时间", ctx.Query("pre_msg_time"), "位数", len(ctx.Query("pre_msg_time")))
	pre_msg_time := ctx.Query("pre_msg_time")
	uid := service.GetIdByToken(ctx.Query("token"))
	toUserId, err := strconv.ParseInt(ctx.Query("to_user_id"), 10, 64)
	if err != nil {
		panic(err)
	}
	var messageLogs []model.Message

	if pre_msg_time == "0" {
		if res := model.Db.Where("origin_user_id=? and destination_user_id=? and create_date > ?", uid, toUserId, pre_msg_time).Or("origin_user_id=? and destination_user_id=? and create_date > ?", toUserId, uid, pre_msg_time).Order("create_date asc").Find(&messageLogs); res.Error != nil {
			panic(res.Error)
		}
	} else {
		if res := model.Db.Where("origin_user_id=? and destination_user_id=? and create_date > ?", toUserId, uid, pre_msg_time).Order("create_date asc").Find(&messageLogs); res.Error != nil {
			panic(res.Error)
		}
	}

	ctx.JSON(200, gin.H{
		"status_code":  "0",
		"status_msg":   "获取成功",
		"message_list": messageLogs,
	})
}

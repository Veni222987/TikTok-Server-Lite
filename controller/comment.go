package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"DoushengABCD/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// Comment 评论操作
func Comment(ctx *gin.Context) {
	token := ctx.Query("token")
	video_id := ctx.Query("video_id")
	action_type := ctx.Query("action_type")
	vid64, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		panic(err)
	}
	uid64 := service.GetIdByToken(token)

	commentLog := model.Comment{UserId: uid64, VideoId: vid64}
	if action_type == "1" {
		//发布评论
		commentText := ctx.Query("comment_text")
		commentLog.Content = commentText
		commentLog.CreateDate = time.Now()
		//更新comment

		if res := service.Db.Create(&commentLog); res.Error != nil {
			panic(res.Error)
		}
		//video表评论数++
		if res := service.Db.Table("video").Where("id=?", vid64).Update("comment_count", gorm.Expr("comment_count+1")); res.Error != nil {
			panic(res.Error)
		}
		var resComment struct {
			Id         int64      `json:"id"`
			User       model.User `json:"user"`
			Content    string     `json:"content"`
			CreateDate string     `json:"create_date"`
		}
		if res := service.Db.Where("id=?", service.GetIdByToken(token)).Find(&resComment.User); res.Error != nil {
			panic(res.Error)
			return
		}
		resComment.Content = commentLog.Content
		resComment.CreateDate = time.Now().Format("2006-01-02 15:04:02")

		ctx.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "评论成功",
			"comment":     resComment,
		})
	} else if action_type == "2" {
		//取消评论
		commentToDel := ctx.Query("comment_id")
		//删除指定评论
		if res := service.Db.Where("comment_id=?", commentToDel).Delete(&commentLog); res.Error != nil {
			panic(res.Error)
		}
		//video表评论数--
		if res := service.Db.Table("video").Where("id=?", vid64).Update("comment_count", gorm.Expr("comment_count-1")); res.Error != nil {
			panic(res.Error)
		}
		ctx.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "取消评论成功",
		})
	} else {
		//非法请求
		ctx.JSON(400, gin.H{
			"status_code": 400,
			"status_msg":  "非法请求",
		})
	}
}

// GetCommentList 获取评论列表
func GetCommentList(ctx *gin.Context) {
	vid64 := ctx.Query("video_id")

	//查找当前视频所有评论的信息
	var commentList []model.Comment
	res := service.Db.Where("video_id=?", vid64).Find(&commentList)
	if res.Error != nil {
		panic(res.Error)
	}

	type resComment struct {
		Id         int              `gorm:"id" json:"id"`
		User       model.User       `json:"user"`
		Content    string           `gorm:"content" json:"content"`
		CreateDate utils.CustomTime `gorm:"create_date" json:"create_date"`
	}

	resCommentList := make([]resComment, len(commentList))
	for i, c := range commentList {
		resCommentList[i] = resComment{
			Id:         c.Id,
			Content:    c.Content,
			CreateDate: utils.CustomTime{c.CreateDate},
		}
		//查找User的详细信息
		res = service.Db.Table("user").Where("id=?", c.UserId).Find(&resCommentList[i].User)
		if res.Error != nil {
			panic(res.Error)
		}
	}

	ctx.JSON(200, gin.H{
		"status_code":  0,
		"status_msg":   "get comment list successfully",
		"comment_list": resCommentList,
	})
}

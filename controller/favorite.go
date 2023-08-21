package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func Like(ctx *gin.Context) {
	video_id := ctx.Query("video_id")
	action_type := ctx.Query("action_type")
	vid64, err := strconv.ParseInt(video_id, 10, 64)
	if err != nil {
		panic(err)
	}
	uid64 := service.GetIdByToken(ctx.Query("token"))
	liekLog := model.Like{UserId: uid64, VideoId: vid64}
	if action_type == "1" {
		//点赞
		//检查是否已经点赞
		var temp []model.Like
		res := model.Db.Where("user_id=? and video_id=?", uid64, vid64).Find(&temp)
		if res.Error != nil {
			panic(res.Error)
		}
		if res.RowsAffected != 0 {
			ctx.JSON(200, gin.H{
				"status_code": 1,
				"status_msg":  "已经点赞过了",
			})
			return
		}
		//封装成为事务，保证数据库的一致性
		tx := model.Db.Begin()
		//更新like
		res = tx.Create(&liekLog)
		fmt.Println("点赞信息", liekLog)
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		//该视频的喜欢数++
		res = tx.Table("video").Where("id=?", vid64).Update("favorite_count", gorm.Expr("favorite_count+1"))
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		//视频作者的喜欢数++
		var authorUid struct {
			AuthorID int64 `gorm:"author_id"`
		}
		res = tx.Table("video").Select("author_id").Where("id=?", vid64).First(&authorUid)
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		res = tx.Table("user").Where("id=?", authorUid.AuthorID).Update("favorite_count", gorm.Expr("favorite_count+1"))
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		//当前用户的喜欢数++
		res = tx.Table("user").Where("id=?", uid64).Update("total_favorited", gorm.Expr("total_favorited+1"))
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		tx.Commit()
		ctx.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "点赞成功",
		})
	} else if action_type == "2" {
		//取消点赞
		//更新like
		tx := model.Db.Begin()
		res := tx.Where("video_id=?", video_id).Delete(&liekLog)
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		//该视频的喜欢数--
		res = tx.Table("video").Where("id=?", vid64).Update("favorite_count", gorm.Expr("favorite_count-1"))
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		//视频作者的喜欢数--
		var authorUid struct {
			AuthorId int64 `gorm:"author_id"`
		}
		res = tx.Table("video").Select("author_id").Where("id=?", vid64).First(&authorUid)
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		res = tx.Table("user").Where("id=?", authorUid.AuthorId).Update("favorite_count", gorm.Expr("favorite_count-1"))
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		//当前用户的喜欢数--
		res = tx.Table("user").Where("id=?", uid64).Update("total_favorited", gorm.Expr("total_favorited-1"))
		if res.Error != nil {
			panic(res.Error)
			tx.Rollback()
			return
		}
		tx.Commit()
		ctx.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "取消点赞成功",
		})
	} else {
		//非法请求
		ctx.JSON(400, gin.H{
			"status_code": 1,
			"status_msg":  "非法请求",
		})
	}
}

// 获取喜欢列表
func GetFavoriteList(ctx *gin.Context) {
	//获取参数
	userID := ctx.Query("user_id")
	//查找Like表中userID的所有视频id
	var likesVID []int64
	res := model.Db.Table("like").Select("video_id").Where("user_id=?", userID).Find(&likesVID)
	if res.Error != nil {
		panic(res.Error)
	}

	//查找视频的所有信息
	type resVideo struct {
		Id            int64      `gorm:"id" json:"id"`
		Author        model.User `json:"author"`
		PlayUrl       string     `gorm:"play_url" json:"play_url"`             // 视频url
		CoverUrl      string     `gorm:"cover_url" json:"cover_url"`           // 封面url
		FavoriteCount int        `gorm:"favorite_count" json:"favorite_count"` // 点赞数量
		CommentCount  int        `gorm:"comment_count" json:"comment_count"`   // 评论数量
		Title         string     `gorm:"title" json:"title"`                   // 视频标题
	}

	var likesVideo []model.Video

	res = model.Db.Table("video").Where("id In ?", likesVID).Find(&likesVideo)
	if res.Error != nil {
		panic(res.Error)
	}
	likesVideoStruct := make([]resVideo, len(likesVideo))
	//对于每个redVideo，找到author信息
	for i, v := range likesVideo {
		likesVideoStruct[i] = resVideo{
			Id:            v.Id,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Title:         v.Title,
		}
		res = model.Db.Where("id =?", v.AuthorId).Find(&likesVideoStruct[i].Author)
		if res.Error != nil {
			panic(res.Error)
		}
	}

	ctx.JSON(200, gin.H{
		"status_code": "0",
		"status_msg":  "成功",
		"video_list":  likesVideoStruct,
	})
}

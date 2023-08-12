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
		//该视频的喜欢数++
		model.Db.Where("id=?", vid64).Update("favorite_count", gorm.Expr("favorite_count+1"))
		//视频作者的喜欢数++
		var authorUid int64
		model.Db.Table("video").Select("author_id").Where("id=?", vid64).First(&authorUid)
		model.Db.Where("id=?", authorUid).Update("favorite_count", gorm.Expr("favorite_count+1"))
		//当前用户的喜欢数++
		model.Db.Where("id=?", uid64).Update("total_favorited", gorm.Expr("total_favorited+1"))
		ctx.JSON(200, gin.H{
			"status_code": 200,
			"status_msg":  "点赞成功",
		})
	} else if action_type == "2" {
		//取消点赞
		//更新like
		model.Db.Where("video_id=?", video_id).Delete(&liekLog)
		//该视频的喜欢数--
		model.Db.Where("id=?", vid64).Update("favorite_count", gorm.Expr("favorite_count-1"))
		//视频作者的喜欢数--
		var authorUid int64
		model.Db.Table("video").Select("author_id").Where("id=?", vid64).First(&authorUid)
		model.Db.Where("id=?", authorUid).Update("favorite_count", gorm.Expr("favorite_count-1"))
		//当前用户的喜欢数--
		model.Db.Where("id=?", uid64).Update("total_favorited", gorm.Expr("total_favorited-1"))
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
		Id            int64 `gorm:"id"`
		Author        model.User
		PlayUrl       string `gorm:"play_url"`       // 视频url
		CoverUrl      string `gorm:"cover_url"`      // 封面url
		FavoriteCount int    `gorm:"favorite_count"` // 点赞数量
		CommentCount  int    `gorm:"comment_count"`  // 评论数量
		Title         string `gorm:"title"`          // 视频标题
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
		"status_code": "200",
		"status_msg":  "成功",
		"video_list":  likesVideoStruct,
	})
}

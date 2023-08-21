package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

func Feed(c *gin.Context) {
	// 获取当前时间戳（秒级）
	currentTime := time.Now().Unix()
	// 获取可选参数
	latestTime := c.Query("latest_time")
	token := c.Query("token")
	var err error
	// 判断时间戳
	if latestTime != "" {
		currentTime, err = strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			fmt.Println("无法将字符串转换为数字", err)
		}
	}
	if token != "" {
		// 获取userName
		userName := service.GetNameByToken(token)
		if userName == "" {
			fmt.Println("key不存在")
		}
	}
	var userID int64
	if token != "" {
		// 获取userID
		userID = service.GetIdByToken(token)
		if userID == 0 {
			fmt.Println("key不存在")
			// 返回数据
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
				"next_time":   math.MaxInt,
				"video_list":  nil,
			})
			return
		}
	}
	// 临时结构体
	// user
	type user struct {
		Avatar          string `json:"avatar"`           // 用户头像
		BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
		FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
		FollowCount     int64  `json:"follow_count"`     // 关注总数
		FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
		ID              int64  `json:"id" gorm:"id"`     // 用户id
		IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
		Name            string `json:"name"`             // 用户名称
		Signature       string `json:"signature"`        // 个人简介
		TotalFavorited  string `json:"total_favorited"`  // 获赞数量
		WorkCount       int64  `json:"work_count"`       // 作品数
	}
	// video
	type video struct {
		AuthorId      int64  `json:"-" gorm:"author_id"`
		Author        user   `json:"author"`                     // 视频作者信息
		CommentCount  int64  `json:"comment_count"`              // 视频的评论总数
		CoverURL      string `json:"cover_url" gorm:"cover_url"` // 视频封面地址
		FavoriteCount int64  `json:"favorite_count"`             // 视频的点赞总数
		ID            int64  `json:"id" gorm:"id"`               // 视频唯一标识
		IsFavorite    bool   `json:"is_favorite"`                // true-已点赞，false-未点赞
		PlayURL       string `json:"play_url" gorm:"play_url"`   // 视频播放地址
		Title         string `json:"title"`                      // 视频标题
		Time          int64  `json:"-"`                          //视频发布时间
		IsFollow      bool   `json:"is_follow"`                  // 是否关注
	}
	var videos []video

	// 查询数据库封装数据
	model.Db.Table("video").Order("time DESC").Limit(30).Where("time <= ?", currentTime).Find(&videos)
	for index, videoT := range videos {
		fmt.Println(videoT.AuthorId)
		var userT user
		model.Db.Table("user").Where("id = ?", videoT.AuthorId).First(&userT)
		videos[index].Author = userT
		var count1 int64
		// 数据库查询是否关注
		model.Db.Table("follow").Where("user_id_a = ? AND user_id_b = ?", userID, userT.ID).Count(&count1)
		if count1 != 0 {
			videos[index].IsFollow = true
		}
		// 数据库查询是否点赞
		var count2 int64
		model.Db.Table("like").Where("user_id = ? AND video_id = ?", userT.ID, videos[index].ID).Count(&count2)
		if count2 != 0 {
			videos[index].IsFavorite = true
		}
	}
	if len(videos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "success",
			"next_time":   9223372036854775807,
			"video_list":  nil,
		})
		return
	}
	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"next_time":   videos[len(videos)-1].Time,
		"video_list":  videos,
	})
	return
}

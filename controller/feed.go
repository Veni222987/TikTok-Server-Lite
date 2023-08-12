package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type FeedResponse struct {
	NextTime   *int64  `json:"next_time"`   // 本次返回的视频中，发布最早的时间，作为下次请求时的latest_time
	StatusCode int64   `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  *string `json:"status_msg"`  // 返回状态描述
	VideoList  []Video `json:"video_list"`  // 视频列表
}

// Video
type Video struct {
	Author        User   `json:"author"`         // 视频作者信息
	CommentCount  int64  `json:"comment_count"`  // 视频的评论总数
	CoverURL      string `json:"cover_url"`      // 视频封面地址
	FavoriteCount int64  `json:"favorite_count"` // 视频的点赞总数
	ID            int64  `json:"id"`             // 视频唯一标识
	IsFavorite    bool   `json:"is_favorite"`    // true-已点赞，false-未点赞
	PlayURL       string `json:"play_url"`       // 视频播放地址
	Title         string `json:"title"`          // 视频标题
	Time          int64  `json:"time"`           //视频发布时间
}

// 视频作者信息

// User
type User struct {
	Avatar          string `json:"avatar"`           // 用户头像
	BackgroundImage string `json:"background_image"` // 用户个人页顶部大图
	FavoriteCount   int64  `json:"favorite_count"`   // 喜欢数
	FollowCount     int64  `json:"follow_count"`     // 关注总数
	FollowerCount   int64  `json:"follower_count"`   // 粉丝总数
	ID              int64  `json:"id"`               // 用户id
	IsFollow        bool   `json:"is_follow"`        // true-已关注，false-未关注
	Name            string `json:"name"`             // 用户名称
	Signature       string `json:"signature"`        // 个人简介
	TotalFavorited  string `json:"total_favorited"`  // 获赞数量
	WorkCount       int64  `json:"work_count"`       // 作品数
}

func Feed(c *gin.Context) {
	// 获取当前时间戳（秒级）
	currentTime := time.Now().Unix()
	// 获取可选参数
	latestTime := c.Query("latest_time")
	token := c.Query("token")
	var err error
	// 根据参数进行相应处理
	if latestTime != "" {
		currentTime, err = strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			// 处理转换错误，例如返回错误信息
			return
		}
	}
	if token != "" {
		// ...
	}

}

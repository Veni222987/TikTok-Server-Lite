package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"DoushengABCD/utils"
	"fmt"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// UploadVideo 投稿接口
func UploadVideo(c *gin.Context) {
	video := model.Video{FavoriteCount: 0, CommentCount: 0}
	// 获取token
	token := c.PostForm("token")
	// 使用 Redis 获取userName
	userName, err := service.RedisClient.Get(token).Result()
	if err == redis.Nil {
		fmt.Println("key不存在")
	} else if err != nil {
		panic("Redistribution，获取userName失败  " + err.Error())
	} else {
		if err != nil {
			c.JSON(1, "userName获取失败")
		}
	}
	// 获取userID
	var user model.User
	model.Db.Table("user").Where("name = ?", userName).First(&user)
	video.AuthorId = user.Id
	// 获取名为 "title"
	video.Title = c.PostForm("title")
	var id int64
	for {
		// 生成视频ID
		id = utils.GenVideoID()
		// 检验唯一性数据库查查询操作
		var count int64
		model.Db.Table("video").Where("id = ?", id).Count(&count)
		if count == 0 {
			video.Id = id
			break
		}
	}
	// c.FormFile函数来获取文件对象
	data, err := c.FormFile("data")
	// 错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 获取文件名，例：bear.mp4
	videoName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%v_%s", user.Id, videoName)
	// 文件名切分
	parts := strings.Split(videoName, ".")
	// 文件名转换
	finalVideoName := strconv.FormatInt(id, 10) + "." + parts[1]
	finalCoverName := strconv.FormatInt(id, 10) + ".jpg"
	video.PlayUrl = "https://abcd-dousheng.oss-cn-guangzhou.aliyuncs.com/videos/" + finalVideoName
	video.CoverUrl = "https://abcd-dousheng.oss-cn-guangzhou.aliyuncs.com/covers/" + finalCoverName
	// 生成视频、封面相对路径
	saveVideoFile := "./public/" + finalVideoName
	saveCoverFile := "./public/" + finalCoverName
	// 保存文件
	if err := c.SaveUploadedFile(data, saveVideoFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 获取封面
	getcover(saveVideoFile, saveCoverFile)
	// 上传视频到阿里云
	err = utils.AliyunOSSUpload("videos", finalVideoName, saveVideoFile)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 上传封面到阿里云
	err = utils.AliyunOSSUpload("covers", finalCoverName, saveCoverFile)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 删除视频文件
	go func() {
		err = os.Remove(saveVideoFile)
		if err != nil {
			println("删除文件失败！！！")
			return
		}
	}()
	// 删除封面文件
	go func() {
		err = os.Remove(saveCoverFile)
		if err != nil {
			println("删除文件失败！！！")
			return
		}
	}()
	// 写入数据库
	video.Time = time.Now().Unix()
	fmt.Println(video)
	result := model.Db.Table("video").Create(&video)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "数据库上传失败",
		})
		return
	}
	// 成功返回响应
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	// 获取参数
	userID := c.Query("user_id")
	// 查询数据库封装数据
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
		AuthorId      int64  `json:"-"`
		Author        user   `json:"author"`                     // 视频作者信息
		CommentCount  int64  `json:"comment_count"`              // 视频的评论总数
		CoverURL      string `json:"cover_url" gorm:"cover_url"` // 视频封面地址
		FavoriteCount int64  `json:"favorite_count"`             // 视频的点赞总数
		ID            int64  `json:"id" gorm:"id"`               // 视频唯一标识
		IsFavorite    bool   `json:"is_favorite"`                // true-已点赞，false-未点赞
		PlayURL       string `json:"play_url" gorm:"play_url"`   // 视频播放地址
		Title         string `json:"title"`                      // 视频标题
		Time          int64  `json:"-"`                          //视频发布时间
	}
	var videos []video
	var user_t user
	model.Db.Table("video").Where("author_id = ?", userID).Find(&videos)
	if len(videos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
			"video_list":  nil,
		})
	}
	for index, video_t := range videos {
		model.Db.Table("user").Where("id = ?", video_t.AuthorId).Find(&user_t)
		videos[index].Author = user_t
		// 数据库查询是否关注

		// 数据库查询是否点赞
		var count int64
		model.Db.Table("like").Where("user_id = ? AND video_id = ?", user_t.ID, videos[index].ID).Count(&count)
		if count != 0 {
			videos[index].IsFavorite = true
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
		"video_list":  videos,
	})
}

// 获取封面
func getcover(videoPath string, coverPath string) error {
	// 执行带环境变量的 ffmpeg 命令
	cmd := exec.Command("./ffmpeg/bin/ffmpeg.exe", "-i", videoPath, "-ss", "00:00:00.000", "-vframes", "1", coverPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

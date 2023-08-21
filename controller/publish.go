package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"DoushengABCD/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

// UploadVideo 投稿接口
func UploadVideo(c *gin.Context) {
	// 初始化
	video := model.Video{FavoriteCount: 0, CommentCount: 0}
	// 获取token
	token := c.PostForm("token")

	video.AuthorId = service.GetIdByToken(token)
	if video.AuthorId == 0 {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "无法查询id",
		})
		return
	}
	// 获取名为 "title"
	video.Title = c.PostForm("title")
	var id int64
	for {
		// 生成视频ID
		id = utils.GenVideoID()
		// 检验唯一性数据库查查询操作
		var count int64
		res := model.Db.Table("video").Where("id = ?", id).Count(&count)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, Response{
				StatusCode: 1,
				StatusMsg:  "数据库查询失败",
			})
			return
		}
		if count == 0 {
			video.Id = id
			break
		}
	}
	// c.FormFile函数来获取文件对象
	data, err := c.FormFile("data")
	// 错误处理
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 3,
			StatusMsg:  "获取文件对象失败 " + err.Error(),
		})
		return
	}
	// 获取文件名，例：bear.mp4
	videoName := filepath.Base(data.Filename)
	//finalName := fmt.Sprintf("%v_%s", video.AuthorId, videoName)
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
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 4,
			StatusMsg:  "文件接收失败 " + err.Error(),
		})
		return
	}
	// 提前返回响应
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
	})

	// 获取封面
	err = getcover(saveVideoFile, saveCoverFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 5,
			StatusMsg:  "封面生成失败 " + err.Error(),
		})
		// 清除视频
		// 删除视频文件
		go func() {
			err = os.Remove(saveVideoFile)
			if err != nil {
				println("删除文件失败！！！")
				return
			}
		}()
		return
	}
	// 上传视频到阿里云
	err = utils.AliyunOSSUpload("videos", finalVideoName, saveVideoFile)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 6,
			StatusMsg:  "视频上传阿里云失败 " + err.Error(),
		})
		return
	}
	// 上传封面到阿里云
	err = utils.AliyunOSSUpload("covers", finalCoverName, saveCoverFile)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 7,
			StatusMsg:  "封面上传阿里云失败 " + err.Error(),
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
	//封装成为事务，保证数据库的一致性
	tx := model.Db.Begin()
	result := tx.Table("video").Create(&video)
	fmt.Println("视频信息", video)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 8,
			StatusMsg:  "数据库上传失败",
		})
		tx.Rollback()
		return
	}
	res := tx.Table("user").Where("id=?", video.AuthorId).Update("work_count", gorm.Expr("work_count+1"))
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 8,
			StatusMsg:  "数据库上传失败2",
		})
		tx.Rollback()
		return
	}
	tx.Commit()
	// 成功返回响应
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "uploaded successfully",
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
		IsFollow      bool   `json:"is_follow"`                  // 是否关注
	}
	var videos []video
	if res := model.Db.Table("video").Where("author_id = ?", userID).Find(&videos); res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 0,
			"status_msg":  "fail",
			"video_list":  nil,
		})
	}
	if len(videos) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 0,
			"status_msg":  "",
			"video_list":  nil,
		})
		return
	}
	for index, videoT := range videos {
		var userT user
		res := model.Db.Table("user").Where("id = ?", videoT.AuthorId).First(&userT)
		videos[index].Author = userT
		if res.Error != nil {
			continue
		}
		// 数据库查询是否关注
		var count1 int64
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
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"video_list":  videos,
	})
}

// 获取封面
func getcover(videoPath string, coverPath string) error {
	// 执行带环境变量的 ffmpeg 命令
	//获取GOOS
	osType := runtime.GOOS
	if osType == "windows" {
		cmd := exec.Command("./utils/ffmpeg.exe", "-i", videoPath, "-ss", "00:00:00.000", "-vframes", "1", coverPath)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	} else if osType == "linux" {
		cmd := exec.Command("ffmpeg", "-i", videoPath, "-ss", "00:00:00.000", "-vframes", "1", coverPath)
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}

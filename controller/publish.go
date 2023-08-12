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
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// UploadVideo 投稿接口
func UploadVideo(c *gin.Context) {
	video := model.Video{FavoriteCount: 0, CommentCount: 0}
	// 从请求的表单中获取名为 "token" 的值，并将其赋值给变量 token。
	token := c.PostForm("token")
	// 获取用户id
	uid, err := service.RedisClient.Get(token).Result()
	if err == redis.Nil {
		fmt.Println("key不存在")
	} else if err != nil {
		panic(err)
	} else {
		video.AuthorId, err = strconv.ParseInt(uid, 10, 64)
		if err != nil {
			c.JSON(1, "用户id获取失败")
		}
	}
	// 从请求的表单中获取名为 "title" 的值。
	video.Title = c.PostForm("title")
	var id int64
	for {
		// 生成视频ID
		id := utils.GenVideoID()
		// 检验唯一性数据库查查询操作
		var count int64
		model.Db.Table("video").Where("id = ?", id).Count(&count)
		if count == 0 {
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
	// 获取文件名例：bear.mp4
	videoName := filepath.Base(data.Filename)
	// 文件名切分
	parts := strings.Split(videoName, ".")
	// 文件名转换
	finalVideoName := string(id) + parts[1]
	finalCoverName := string(id) + ".jpg"
	video.PlayUrl = "https://oss-cn-guangzhou.aliyuncs.com/videos/" + finalVideoName
	video.CoverUrl = "https://oss-cn-guangzhou.aliyuncs.com/covers/" + finalCoverName
	// 生成视频、封面相对路径
	saveVideoFile := filepath.Join("./public/", finalVideoName)
	saveCoverFile := filepath.Join("./public", finalCoverName)
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
	result := model.Db.Table("video").Create(&video)
	if result.Error != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 2,
			StatusMsg:  "数据库上传失败",
		})
	}
	// 成功返回响应
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  " uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	// 获取参数
	token := c.Query("token")
	// 验证用户身份
	println(token)
	user_id := c.Query("user_id")
	// 查询数据库封装数据
	println(user_id)
}

// 获取封面
func getcover(videoPath string, coverPath string) error {
	// 执行带环境变量的 ffmpeg 命令
	cmd := exec.Command("D:\\install\\ffmpeg-6.0-essentials_build\\bin\\ffmpeg.exe", "-i", videoPath, "-ss", "00:00:00.000", "-vframes", "1", coverPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

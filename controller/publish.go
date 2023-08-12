package controller

import (
	"DoushengABCD/utils"
	_ "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 投稿接口
func UploadVideo(c *gin.Context) {
	//1 验证
	// 从请求的表单中获取名为 "token" 的值，并将其赋值给变量 token。
	token := c.PostForm("token")
	// 检查用户的token是否有效，数据库操作,从token中解析出用户信息
	//2 获取信息
	println(token)
	// 从请求的表单中获取名为 "title" 的值，并将其赋值给变量 title。
	title := c.PostForm("title")
	// 生成视频ID
	id := GenID()
	// 检验唯一性数据库查查询操作

	// 写入数据库
	println(title)
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
	// 数据存入数据库
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
	// 成功返回响应
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fileName + " uploaded successfully",
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

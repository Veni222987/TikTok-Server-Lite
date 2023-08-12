package controller

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 初始化
func init() {
	err := Init("2023-08-11", 1)
	if err != nil {
		println("初始化失败！", err)
	}
}

var node *snowflake.Node

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
	err = aliyunOSS("videos", finalVideoName, saveVideoFile)
	// 错误处理
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	// 上传封面到阿里云
	err = aliyunOSS("covers", finalCoverName, saveCoverFile)
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

// 雪花算法生成视频id
// 初始化雪花算法
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err // 返回错误信息
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Println(err)
		return err // 返回错误信息
	}
	return nil // 返回 nil 表示初始化成功
}

// 生成 64 位的雪花 ID
func GenID() int64 {
	return node.Generate().Int64()
}

// 存储视频到阿里云oss
func aliyunOSS(t string, fileName string, localFilePath string) error {
	// yourBucketName填写存储空间名称。
	bucketName := "abcd-dousheng"
	// yourObjectName填写Object完整路径，完整路径不包含Bucket名称。
	var objectName string
	switch t {
	case "videos":
		objectName = "videos/" + fileName
	case "covers":
		objectName = "covers/" + fileName
	}
	// yourLocalFileName填写本地文件的完整路径或相对路径。
	localFileName := localFilePath
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	client, err := oss.New("https://oss-cn-guangzhou.aliyuncs.com", "LTAI5t9K833uebsoQvSDZxDH", "VO2AUub801jTf6KFczQZoRf4CZhJY0")
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return err
	}
	return nil
}

// 获取封面
func getcover(videoPath string, coverPath string) error {
	// 执行带环境变量的 ffmpeg 命令
	cmd := exec.Command("D:\\install\\ffmpeg-6.0-essentials_build\\bin\\ffmpeg.exe", "-i", videoPath, "-ss", "00:00:00.000", "-vframes", "1", coverPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
}

// 封装视频信息

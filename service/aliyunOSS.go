package service

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"gopkg.in/yaml.v3"
	"os"
)

type OSSConfig struct {
	Endpoint     string `yaml:"endPoint"`
	AccessKey    string `yaml:"accessKey"`
	AccessSecret string `yaml:"accessSecret"`
}

/*
上传文件到阿里云OOS,访问路径：“https://abcd-dousheng.oss-cn-guangzhou.aliyuncs.com/文件类型/文件名”
t: 上传文件类型（“videos”:视频，“covers”：视频封面，“avatar”：用户封面）
fileName：文件名（例如xxx.mp4）
localFilePath:本地文件路径（相对根目录路径或绝对路径）
return：报错
*/
func AliyunOSSUpload(t string, fileName string, localFilePath string) error {
	var ossConfig struct {
		OSS OSSConfig `yaml:"OSS"`
	}
	//读取阿里云OSS配置
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configFile, &ossConfig)
	if err != nil {
		panic(err)
	}
	// yourBucketName填写存储空间名称。
	bucketName := "abcd-dousheng"
	// yourObjectName填写Object完整路径，完整路径不包含Bucket名称。
	var objectName string
	switch t {
	case "videos":
		objectName = "videos/" + fileName
	case "covers":
		objectName = "covers/" + fileName
	case "avatar":
		objectName = "avatar/" + fileName
	default:
		return errors.New("第一个参数不存在！")
	}
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint
	client, err := oss.New(ossConfig.OSS.Endpoint, ossConfig.OSS.AccessKey, ossConfig.OSS.AccessSecret)
	//fmt.Printf("#+v", ossConfig)
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return err
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFilePath)
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
}

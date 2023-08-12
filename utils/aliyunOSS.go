package utils

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

/*
上传文件到阿里云OOS,访问路径：“https://abcd-dousheng.oss-cn-guangzhou.aliyuncs.com/文件类型/文件名”
t: 上传文件类型（“videos”:视频，“covers”：视频封面，“avatar”：用户封面）
fileName：文件名（例如xxx.mp4）
localFilePath:本地文件路径（相对根目录路径或绝对路径）
return：报错
*/
func AliyunOSSUpload(t string, fileName string, localFilePath string) error {
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
	// yourLocalFileName填写本地文件的完整路径或相对路径。
	localFileName := localFilePath
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint
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

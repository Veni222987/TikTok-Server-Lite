package utils

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("asdfghjklzxcvbnmqwertyuiopASDFGHJKLZXCVBNMQWERTYUIOP")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func GenerateToken(username string, uid int64) string {
	//使用用户名+udi生成token
	str := strconv.FormatInt(uid, 10)
	return username + str
}

//var SnowflakeNode snowflake.Node

var videoNode *snowflake.Node
var userIDNode *snowflake.Node

// 雪花算法生成视频id
// 初始化雪花算法
func initHelp(startTime string, machineID int64, node **snowflake.Node) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err // 返回错误信息
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	*node, err = snowflake.NewNode(machineID)
	if err != nil {
		fmt.Println(err)
		return err // 返回错误信息
	}
	return nil // 返回 nil 表示初始化成功
}

// 初始化
func init() {
	res := strings.Split(time.Now().String(), " ")
	err := initHelp(res[0], 1, &videoNode)
	if err != nil {
		println("初始化失败！", err)
	}

	err = initHelp(res[0], 1, &userIDNode)
	if err != nil {
		println("初始化失败", err)
	}
}

// 生成 64 位的视频ID
func GenVideoID() int64 {
	return videoNode.Generate().Int64()
}

// 生成User_id
func GenUserID() int64 {
	return userIDNode.Generate().Int64()
}

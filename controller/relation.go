package controller

import (
	"DoushengABCD/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 关注操作
func RelationAction(c *gin.Context) {
	toUserID := c.Query("to_user_id")
	actionType := c.Query("action_type")
	// 从Redis获取自己的ID
	fromUserID := service.GetIdByToken(c.Query("token"))
	// 写入数据库，错误处理

	// 成功情况返回
	c.JSON(http.StatusOK, c.H{
		"status_code": 0,
		"status_msg":  "success",
	})
}

// 关注列表
func FollowList(c *gin.Context) {
	userID := c.Query("user_id")
	// 查询数据库

	//返回

}

// 粉丝列表
func FollowerList(c *gin.Context) {
	userID := c.Query("user_id")
	// 查询数据库

	// 封装返回结果

	// 返回

}

// 好友列表
func FriendList(c *gin.Context) {
	userID := c.Query("user_id")
	// 查询数据库

	// 封装返回结果

	// 返回

}

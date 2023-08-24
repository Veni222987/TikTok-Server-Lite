package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// RelationAction 关注与取消关注操作
func RelationAction(c *gin.Context) {
	// 获取参数
	toUserID, err := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 1,
			"status_msg":  "fail",
		})
		return
	}
	actionType := c.Query("action_type")
	// 从Redis获取自己的ID
	fromUserID := service.GetIdByToken(c.Query("token"))
	follow := model.Follow{UserIdA: fromUserID, UserIdB: toUserID}
	//封装成为事务，保证数据库的一致性
	tx := service.Db.Begin()
	if actionType == "1" {
		// 关注
		// follow表
		if follow.UserIdA == follow.UserIdB {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "无法关注自己",
			})
			return
		}
		var count int64
		res := tx.Table("follow").Where("user_id_a = ? AND user_id_b = ?", follow.UserIdA, follow.UserIdB).Count(&count)
		if res.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "数据库查询错误",
			})
			// 回滚
			tx.Rollback()
			return
		}
		if count == 1 {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "无法重复关注",
			})
			// 回滚
			tx.Rollback()
			return
		}
		res = tx.Table("follow").Create(&follow)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
			})
			// 回滚
			tx.Rollback()
			return
		}
		// 关注者，关注数++
		res = tx.Table("user").Where("id=?", fromUserID).Update("follow_count", gorm.Expr("follow_count+1"))
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
			})
			// 回滚
			tx.Rollback()
			return
		}
		// 被关注者，粉丝数++
		res = tx.Table("user").Where("id=?", toUserID).Update("follower_count", gorm.Expr("follower_count+1"))
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
			})
			// 回滚
			tx.Rollback()
			return
		}
	} else if actionType == "2" {
		// 取消关注
		if follow.UserIdA == follow.UserIdB {
			c.JSON(http.StatusOK, gin.H{
				"status_code": 1,
				"status_msg":  "无法关注自己",
			})
			return
		}
		//封装成为事务，保证数据库的一致性
		res := service.Db.Table("follow").Where("user_id_a = ? AND user_id_b = ?", fromUserID, toUserID).Delete(&follow)
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
			})
			// 回滚
			tx.Rollback()
			return
		}
		// 关注者，关注数--
		res = tx.Table("user").Where("id=?", fromUserID).Update("follow_count", gorm.Expr("follow_count-1"))
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
			})
			// 回滚
			tx.Rollback()
			return
		}
		// 被关注者，粉丝数--
		res = tx.Table("user").Where("id=?", toUserID).Update("follower_count", gorm.Expr("follower_count-1"))
		if res.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "fail",
			})
			// 回滚
			tx.Rollback()
			return
		}

	}
	// 成功情况返回
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "success",
	})
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
	userID := c.Query("user_id")
	type user struct {
		model.User
		IsFollow bool `json:"is_follow"`
	}
	var followList []model.Follow
	var userList []user
	// 查询数据库
	res := service.Db.Table("follow").Where("user_id_a = ?", userID).Find(&followList)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": "1",
			"status_msg":  "fail",
			"user_list":   nil,
		})
		return
	}
	for _, temp := range followList {
		var userT user
		res := service.Db.Table("user").Where("id = ?", temp.UserIdB).Find(&userT)
		if res.Error != nil {
			continue
		}
		userT.IsFollow = true
		userList = append(userList, userT)
	}
	//返回
	c.JSON(http.StatusOK, gin.H{
		"status_code": "0",
		"status_msg":  "success",
		"user_list":   userList,
	})

}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	userID := c.Query("user_id")
	type user struct {
		model.User
		IsFollow bool `json:"is_follow"`
	}
	var followList []model.Follow
	var userList []user
	// 查询数据库
	// 找谁关注我
	res := service.Db.Table("follow").Where("user_id_b = ?", userID).Find(&followList)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": "1",
			"status_msg":  "fail",
			"user_list":   nil,
		})
		return
	}
	for _, temp := range followList {
		var userT user
		res := service.Db.Table("user").Where("id = ?", temp.UserIdA).Find(&userT)
		if res.Error != nil {
			continue
		}
		// 判断我是否关注了
		var count int64
		res = service.Db.Table("follow").Where("user_id_a = ? AND user_id_b = ?", userID, userT.Id).Count(&count)
		if res.Error != nil {
			continue
		}
		if count != 0 {
			userT.IsFollow = true
		}
		userList = append(userList, userT)
	}
	//返回
	c.JSON(http.StatusOK, gin.H{
		"status_code": "0",
		"status_msg":  "success",
		"user_list":   userList,
	})

}

// FriendList 好友列表
func FriendList(c *gin.Context) {
	userID := c.Query("user_id")
	var followList []model.Follow
	var userList []model.User
	// 查询数据库
	// 我关注的谁
	res := service.Db.Table("follow").Where("user_id_a = ?", userID).Find(&followList)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": "1",
			"status_msg":  "fail",
			"user_list":   nil,
		})
		return
	}
	for _, temp := range followList {
		var userT model.User
		res := service.Db.Table("user").Where("id = ?", temp.UserIdB).Find(&userT)
		if res.Error != nil {
			continue
		}
		// 它是否有关注我
		var count int64
		res = service.Db.Table("follow").Where("user_id_a = ? AND user_id_b = ?", userT.Id, userID).Count(&count)
		if res.Error != nil {
			continue
		}
		if count == 0 {
			continue
		}
		userList = append(userList, userT)
	}
	//返回
	c.JSON(http.StatusOK, gin.H{
		"status_code": "0",
		"status_msg":  "success",
		"user_list":   userList,
	})

}

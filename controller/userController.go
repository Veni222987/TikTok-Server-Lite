package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {
	DB := model.Db
	//获取参数
	name := ctx.Query("username") //注意字符串要用双引号
	//telephone := ctx.Query("telephone")
	password := ctx.Query("password")
	////数据验证
	//if len(telephone) != 11 {
	//	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
	//	return
	//}
	//if len(password) < 6 {
	//	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
	//	return
	//}

	//如果没有传name（名称），就给他一个10位的随机字符串
	//if len(name) == 0 {
	//	name = utils.RandomString(10)
	//}
	//log.Println(name, telephone, password)

	//判断手机号是否存在
	//if isUserExist(DB, telephone) {
	//	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
	//}
	//如果用户不存在，则创建用户
	if isUserExist(DB, name) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
	}

	//对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "加密错误"})
		return
	}

	newAccount := model.Account{name, string(hashedPassword)}

	res := DB.Create(&newAccount)
	if res.Error != nil {
		log.Println(res.Error)
	}

	//创建用户
	id := utils.GenerateId(name)
	user := model.User{Id: id, Name: name}
	res = DB.Create(&user)
	if res.Error != nil {
		log.Println(res.Error)
	}

	//生成token并保存到redis
	token := utils.GenerateToken(name)
	log.Println(token)

	//返回结果
	ctx.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "注册成功",
		"user_id":     id,
		"token":       utils.GenerateToken(name),
	})
}

// 登录功能
func Login(ctx *gin.Context) {
	DB := model.Db
	//获取参数
	name := ctx.PostForm("username")
	password := ctx.PostForm("password")
	//数据验证
	//if len(telephone) != 11 {
	//	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
	//	return
	//}
	//if len(password) < 6 {
	//	ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
	//	return
	//}
	//判断用户是否存在
	account := model.Account{}
	DB.Where("username = ?", name).First(&account)
	if len(account.Username) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}
	//发送token
	token := utils.GenerateToken(name)
	//返回结果
	ctx.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "登陆成功",
		"user_id":     0,
		"token":       token,
	})
}

// 判断手机号是否存在
func isUserExist(db *gorm.DB, username string) bool {
	var account model.Account
	model.Db.Select("username").Where("username =?", username).First(&account)
	if len(account.Username) != 0 {
		return true
	}
	return false
}

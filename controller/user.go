package controller

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"DoushengABCD/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// Register 用户注册
func Register(ctx *gin.Context) {
	DB := service.Db
	//获取参数
	name := ctx.Query("username") //注意字符串要用双引号
	password := ctx.Query("password")

	//判断用户是否存在，如果用户不存在，则创建用户
	if isUserExist(DB, name) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
		return
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

	//生成User_id，创建用户
	id := utils.GenUserID()
	user := model.User{Id: id, Name: name}
	res = DB.Omit("avatar", "background_image", "signature").Create(&user)
	if res.Error != nil {
		log.Println(res.Error)
	}

	//生成token并保存到redis，过期时间为1天
	token := utils.GenerateToken(name, id)
	userInfo := map[string]interface{}{
		"id":   id,
		"name": name,
	}
	userInfoJson, _ := json.Marshal(userInfo)
	service.RedisClient.Set(token, userInfoJson, 86400000000000)
	//log.Println(token)

	//返回结果
	ctx.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "注册成功",
		"user_id":     id,
		"token":       token,
	})
}

// Login 登录功能
func Login(ctx *gin.Context) {
	DB := service.Db
	//获取参数
	name := ctx.Query("username")
	password := ctx.Query("password")

	//判断用户是否存在
	account := model.Account{}
	DB.Table("account").Where("username = ?", name).Find(&account)
	if len(account.Username) == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//获取user_id
	var uid64 int64

	res := service.Db.Table("user").Select("id").Where("name=?", name).Find(&uid64)
	if res.Error != nil {
		panic(res.Error)
	}
	//发送token，过期时间为1天
	token := utils.GenerateToken(name, uid64)

	type uInfoStruct struct {
		ID   int64  `gorm:"id" json:"id"`
		Name string `gorm:"name" json:"name"`
	}
	var uInfo uInfoStruct

	service.Db.Table("user").Select("id,name").Where("name=?", name).First(&uInfo)
	//序列化
	userInfoJson, err := json.Marshal(uInfo)

	err = service.RedisClient.Set(token, userInfoJson, 86400000000000).Err()
	if err != nil {
		panic(err)
	}

	//返回结果
	ctx.JSON(200, gin.H{
		"status_code": 0,
		"status_msg":  "string",
		"user_id":     uid64,
		"token":       token,
	})
}

// 判断用户名是否存在
func isUserExist(db *gorm.DB, username string) bool {
	var account model.Account
	service.Db.Select("username").Where("username =?", username).First(&account)
	if len(account.Username) != 0 {
		return true
	}
	return false
}

func GetUserInfo(ctx *gin.Context) {
	//获取参数
	userID := ctx.Query("user_id")
	//fmt.Println("接收id", userID)
	//token := ctx.Query("token")
	//根据userID查找数据库
	var user model.User
	if res := service.Db.Where("id=?", userID).Find(&user); res.Error != nil {
		panic(res.Error)
	}

	type resUser struct {
		Id            int
		Name          string
		FollowCount   int
		FollowerCount int
		IsFollow      bool
	}

	ctx.JSON(0, gin.H{
		"status_code": 0,
		"status_msg":  "success",
		"user":        user,
	})
}

package main

import (
	"DoushengABCD/model"
	"DoushengABCD/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 获取结构体
	//utils.TableConverter()

	r := gin.Default()

	//初始化路由
	InitRouter(r)
	model.InitDatabases()
	service.InitRedis()
	
	//u := model.User{Id: 111, Name: "Veni"}
	//
	//res := model.Db.Create(&u)
	//if err := res.Error; err != nil {
	//	fmt.Println("插入失败", err)
	//	return
	//}

	r.Run(":8888")
}

package main

import (
	"DoushengABCD/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//初始化路由
	InitRouter(r)
	InitDatabases()

	u := model.User{111, "Veni"}

	res := Db.Create(&u)
	if err := res.Error; err != nil {
		fmt.Println("插入失败", err)
		return
	}

	r.Run(":8888")
}

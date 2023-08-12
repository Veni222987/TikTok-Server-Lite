package main

import (
	"ginEssential/common"
	"github.com/gin-gonic/gin"
)

func main() {
	db := common.InitDB()
	//获取底层数据库连接
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database,err: " + err.Error())
	}
	defer sqlDB.Close()
	r := gin.Default()
	r = collectRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080
}

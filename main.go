package main

import (
	"DoushengABCD/service"
	"github.com/gin-gonic/gin"
)

var Test string

func main() {
	r := gin.Default()
	//初始化
	InitRouter(r)
	service.InitDatabase()
	service.InitRedis()

	r.Run(":8888")
}

package utils

// 参考：https://github.com/gohouse/converter
import (
	"fmt"
	"github.com/gohouse/converter"
)

func TableConverter() {
	err := converter.NewTable2Struct().
		SavePath("./model/model.go").
		TagKey("gorm").
		// 是否添加结构体方法获取表名
		RealNameMethod("TableName").
		Dsn("username:password@tcp(127.0.0.1:3306)/dousheng?charset=utf8").
		Run()
	fmt.Println(err)
}

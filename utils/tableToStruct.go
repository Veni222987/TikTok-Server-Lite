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
		Dsn("hu1hu:asdfghjkl@tcp(47.113.149.158:3306)/dousheng?charset=utf8").
		Run()
	fmt.Println(err)
}

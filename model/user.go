package model

type User struct {
	Id   int    `gorm:"column:id"`
	Name string `gorm:"column:name"`
}

func (u User) TableName() string {
	return "user"
}

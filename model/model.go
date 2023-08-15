package model

import "time"

type User struct {
	Id              int64  `gorm:"id" json:"id"`                             // 用户id
	Name            string `gorm:"name" json:"name"`                         // 用户名称
	FollowCount     int    `gorm:"follow_count" json:"follow_count"`         // 关注总数
	FollowerCount   int    `gorm:"follower_count" json:"follower_count"`     // 粉丝总数
	Avatar          string `gorm:"avatar" json:"avatar"`                     // 用户头像
	BackgroundImage string `gorm:"background_image" json:"background_image"` // 用户个人页顶部大图
	Signature       string `gorm:"signature" json:"signature"`               // 个人简介
	TotalFavorited  int    `gorm:"total_favorited" json:"total_favorited"`   // 获赞数量
	WorkCount       int    `gorm:"work_count" json:"work_count"`             // 作品数
	FavoriteCount   int    `gorm:"favorite_count" json:"favorite_count"`     // 喜欢数
}

func (*User) TableName() string {
	return "user"
}

type Account struct {
	Username string `gorm:"username"`
	Password string `gorm:"password"`
}

func (*Account) TableName() string {
	return "account"
}

type Video struct {
	Id            int64  `gorm:"id"`
	AuthorId      int64  `gorm:"author_id"`      // 作者id
	PlayUrl       string `gorm:"play_url"`       // 视频url
	CoverUrl      string `gorm:"cover_url"`      // 封面url
	FavoriteCount int    `gorm:"favorite_count"` // 点赞数量
	CommentCount  int    `gorm:"comment_count"`  // 评论数量
	Title         string `gorm:"title"`          // 视频标题
	Time          int64  `gorm:"time"`           // 时间戳
}

func (*Video) TableName() string {
	return "video"
}

type Comment struct {
	Id         int       `gorm:"id"`
	UserId     int64     `gorm:"user_id"`
	VideoId    int64     `gorm:"video_id"`
	Content    string    `gorm:"content"`
	CreateDate time.Time `gorm:"create_date"`
}

func (*Comment) TableName() string {
	return "comment"
}

type Like struct {
	Id      int   `gorm:"id"`
	UserId  int64 `gorm:"user_id"`
	VideoId int64 `gorm:"video_id"`
}

func (*Like) TableName() string {
	return "like"
}

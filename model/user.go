package model

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Work struct {
	Title    string    `json:"title"`    //作品名
	Image    string    `json:"image"`    //链接
	Content  string    `json:"content"`  //文字介绍
	Author   string    `json:"author"`   //作者名
	Comments []Comment `json:"comments"` //评论
}
type Comment struct {
	FromUser  string    `json:"from_user"`  //评论者
	Content   string    `json:"content"`    //评论
	CreatedAt time.Time `json:"created_at"` //评论时间
}

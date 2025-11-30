package model

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Work struct {
	Title   string `json:"title"`   //作品名
	Image   string `json:"image"`   //链接
	Content string `json:"content"` //文字介绍
	Author  string `json:"author"`  //作者名
}

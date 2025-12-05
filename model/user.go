package model

import "time"

//用户
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//作品
type Work struct {
	Title    string    `json:"title"`    //作品名
	Image    string    `json:"image"`    //链接
	Content  string    `json:"content"`  //文字介绍
	Author   string    `json:"author"`   //作者名
	Comments []Comment `json:"comments"` //评论
}

//评论
type Comment struct {
	FromUser  string    `json:"from_user"`  //评论者
	Content   string    `json:"content"`    //评论
	CreatedAt time.Time `json:"created_at"` //评论时间
}

// 评论要求（发表评论用）
type CommentRequest struct {
	TargetAuthor string `json:"target_author"` // 评论谁的画
	WorkTitle    string `json:"work_title"`    // 哪幅画
	Content      string `json:"content"`       // 评论内容
}

// [新增] 删除评论的专用请求包
// 为了防止 Body 读取一次就失效，我们需要把所有参数打包在一个结构体里一次性传进来
type DeleteCommentReq struct {
	Owner     string    `json:"owner"`      // 这幅画挂在谁的画廊里 (画展主人)
	Title     string    `json:"title"`      // 作品标题
	FromUser  string    `json:"from_user"`  // 这条评论是谁写的
	CreatedAt time.Time `json:"created_at"` // 这条评论是啥时候写的 (用于精确匹配)
	Content   string    `json:"content"`    // 评论内容 (双重保险匹配)
}

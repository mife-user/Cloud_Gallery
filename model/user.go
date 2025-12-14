package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户
type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(100)" json:"username"`
	Password string `gorm:"column:password;type:varchar(255)" json:"password"`
	UserHand string `gorm:"colume:userhand;type:varchar(500)" json:"userhand"`
	Work     []Work `gorm:"foreignKey:UserID"`
}

// 作品
type Work struct {
	gorm.Model
	UserID   uint      `gorm:"column:user_id"`
	Title    string    `gorm:"column:title;type:varchar(255)" form:"title" json:"title"`
	Image    string    `gorm:"column:image;type:varchar(500)" form:"image" json:"image"`
	Content  string    `gorm:"column:content;type:longtext" form:"content" json:"content"`
	Author   string    `gorm:"column:author;type:varchar(100)" form:"author" json:"author"`
	Comments []Comment `gorm:"foreignKey:WorkID" json:"comments"`
}

// 评论
type Comment struct {
	gorm.Model
	WorkID    uint      `gorm:"column:work_id"`
	FromUser  string    `gorm:"column:fromuser;type:varchar(100)" json:"from_user"`
	Content   string    `gorm:"column:content;type:longtext" json:"content"`
	CreatedAt time.Time `json:"created_at"`
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

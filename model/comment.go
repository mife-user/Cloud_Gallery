package model

import (
	"time"

	"gorm.io/gorm"
)

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
	TargetAuthor string `gorm:"column:target_author;type:varchar(100)" json:"target_author"` // 评论谁的画
	WorkTitle    string `gorm:"column:work_title;type:varchar(255)" json:"work_title"`       // 哪幅画
	Content      string `gorm:"column:content;type:longtext" json:"content"`                 // 评论内容
}

// 删除评论的专用请求包
type DeleteCommentReq struct {
	Owner     string    `json:"owner"`      // 这幅画挂在谁的画廊里 (画展主人)
	Title     string    `json:"title"`      // 作品标题
	FromUser  string    `json:"from_user"`  // 这条评论是谁写的
	CreatedAt time.Time `json:"created_at"` // 这条评论是啥时候写的 (用于精确匹配)
	Content   string    `json:"content"`    // 评论内容 (双重保险匹配)
}

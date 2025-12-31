package model

import (
	"time"

	"gorm.io/gorm"
)

// 评论
type Comment struct {
	gorm.Model
	WorkID   uint   `gorm:"column:workid;type:int"`
	FromUser string `gorm:"column:fromuser;type:varchar(100)" json:"fromuser"`
	Content  string `gorm:"column:content;type:longtext" json:"content"`
}

// 评论要求（发表评论用）
type CommentRequest struct {
	TargetAuthor string `gorm:"column:targetauthor;type:varchar(100)" json:"targetauthor"` // 评论谁的画
	WorkTitle    string `gorm:"column:worktitle;type:varchar(255)" json:"worktitle"`       // 哪幅画
	Content      string `gorm:"column:content;type:longtext" json:"content"`               // 评论内容
}

// 删除评论的专用请求包
type DeleteCommentReq struct {
	Owner     string    `gorm:"column:owner;type:varchar(100)" form:"owner" json:"owner"`          // 这幅画挂在谁的画廊里 (画展主人)
	Title     string    `gorm:"column:title;type:varchar(255)" form:"title" json:"title"`          // 作品标题
	FromUser  string    `gorm:"column:fromuser;type:varchar(100)" form:"fromuser" json:"fromuser"` // 这条评论是谁写的
	CreatedAt time.Time `gorm:"column:createdat" form:"createdat" json:"createdat"`                // 这条评论是啥时候写的 (用于精确匹配)
	Content   string    `gorm:"column:content;type:longtext" form:"content" json:"content"`        // 评论内容 (双重保险匹配)
}

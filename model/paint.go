package model

import "gorm.io/gorm"

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

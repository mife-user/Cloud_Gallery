package model

import (
	"gorm.io/gorm"
)

// 用户
type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(100)" form:"username" json:"username"`
	Password string `gorm:"column:password;type:varchar(255)" form:"password" json:"password"`
	UserHand string `gorm:"colume:userhand;type:varchar(500)" form:"userhand" json:"userhand"`
	Work     []Work `gorm:"foreignKey:userid"`
}

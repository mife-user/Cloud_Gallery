package model

import (
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

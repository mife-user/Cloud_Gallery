package api

import (
	"painting/model"
	"painting/web/service"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	service.UserRegister(c, &u)
}

// 登录
func Login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	service.UserLogin(c, &u)
}

// 添加头像
func AddUserHand(c *gin.Context) {
	service.UserHand(c)
}

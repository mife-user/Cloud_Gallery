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
	username, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "未登录"})
		return
	}
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	user.Username = username.(string)
	service.UserHand(c, &user)
}

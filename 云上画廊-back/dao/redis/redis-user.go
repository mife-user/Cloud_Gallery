package redis

import (
	"painting/box"
	"painting/model"
	"painting/web/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册缓存检查
func Register_read(c *gin.Context, u *model.User) bool {
	if err := box.Temp.RE.HGet(c, u.Username, "username").Err(); err == nil {
		c.JSON(400, gin.H{"error": "用户已存在"})
		return false
	}
	return true
}

// 注册缓存设置
func Register_write(c *gin.Context, u *model.User) {
	if err := box.Temp.RE.HSet(c, u.Username, "username", u.Username).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
	box.Temp.RE.Expire(c, u.Username, 2*time.Hour)
}

// 登录缓存检查与处理
func Login_read(c *gin.Context, u *model.User) int {
	userTemp, err := box.Temp.RE.HMGet(c, u.Username, "username", "password").Result()
	if err != nil {
		return 0
	}
	if userTemp[0] == nil || userTemp[1] == nil {
		return 0
	}
	userName, ok1 := userTemp[0].(string)
	passWord, ok2 := userTemp[1].(string)
	if !ok1 || !ok2 {
		return 0
	}
	if ok := box.Temp.CheckUser(userName, passWord); ok {
		token, _ := utils.GenerateToken(userName)
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		return 2
	} else {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return 1
	}
}

// 登录缓存设置与处理
func Login_write(c *gin.Context, u *model.User) {
	if err := box.Temp.RE.HMSet(c, u.Username,
		"username", u.Username,
		"password", u.Password,
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
	box.Temp.RE.Expire(c, u.Username, 2*time.Hour)
}

// 添加头像缓存处理
func AddUserHand_write(c *gin.Context, user *model.User) {
	if err := box.Temp.RE.HSet(c, user.Username, "userhand", user.UserHand).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

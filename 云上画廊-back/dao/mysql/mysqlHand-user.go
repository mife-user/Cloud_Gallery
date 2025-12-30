package mysql

import (
	"fmt"
	"os"
	"painting/box"
	"painting/dao/redis"
	"painting/model"
	"painting/web/utils"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 数据库注册
func Register(c *gin.Context, u *model.User) int {
	if ok := box.Temp.AddUser(u.Username, u.Password); ok {
		c.JSON(200, gin.H{"message": "注册成功！"})
		return 1
	}
	return 2
}

// 数据库登录
func Login(c *gin.Context, u *model.User) int {
	if ok := box.Temp.CheckUser(u.Username, u.Password); ok {
		token, _ := utils.GenerateToken(u.Username)
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		return 1
	} else {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return 2
	}
}

// 数据库添加头像
func AddUserHand(c *gin.Context, user *model.User) {
	if user.Username == "" {
		c.JSON(400, gin.H{"error": "username required"})
		return
	}
	file, err := c.FormFile("userhand")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}
	dstName := fmt.Sprintf("%s_%d_%s", user.Username, time.Now().Unix(), filepath.Ext(file.Filename))
	if err := os.MkdirAll("userhands", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return
	}
	filePath := filepath.Join("userhands", dstName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return
	}
	c.JSON(200, gin.H{"message": "ok", "path": "/userhands/" + dstName})
	user.UserHand = "/userhands/" + dstName
	box.Temp.AddHand(user)
	redis.AddUserHand_write(c, user)
}

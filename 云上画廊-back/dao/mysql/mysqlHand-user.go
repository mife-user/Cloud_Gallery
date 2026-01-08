package mysql

import (
	"errors"
	"fmt"
	"os"
	"painting/box"
	"painting/model"
	"painting/web/utils"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 数据库注册
func Register(c *gin.Context, u *model.User) error {
	if err := box.Temp.AddUser(u.Username, u.Password); err != nil {
		c.JSON(400, gin.H{"error": "数据库添加失败"})
		return err
	} else {
		c.JSON(200, gin.H{"message": "注册成功！"})
		return nil
	}

}

// 数据库登录
func Login(c *gin.Context, u *model.User) error {
	if err := box.Temp.CheckUser(u.Username, u.Password); err != nil {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return err
	} else {
		token, errToken := utils.GenerateToken(u.Username)
		if errToken != nil {
			c.JSON(400, gin.H{"error": "密码解析错误"})
			return errToken
		}
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		return nil
	}
}

// 数据库添加头像
func AddUserHand(c *gin.Context, user *model.User) error {
	if user.Username == "" {
		c.JSON(400, gin.H{"error": "username required"})
		return errors.New("无名")
	}
	file, err := c.FormFile("userhand")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return errors.New("文件缺失")
	}
	dstName := fmt.Sprintf("%s_%d_%s", user.Username, time.Now().Unix(), filepath.Ext(file.Filename))
	if err := os.MkdirAll("userhands", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return errors.New("头像目录创建失败")
	}
	filePath := filepath.Join("userhands", dstName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return errors.New("保存文件失败")
	}
	c.JSON(200, gin.H{"message": "ok", "path": "/userhands/" + dstName})
	user.UserHand = "/userhands/" + dstName
	box.Temp.AddHand(user)
	return nil
}

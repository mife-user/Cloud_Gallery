package redis

import (
	"errors"
	"painting/box"
	"painting/model"
	"painting/web/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册缓存检查
func Register_read(c *gin.Context, u *model.User) error {
	if err := box.Temp.RE.HGet(c, u.Username, "username").Err(); err == nil {
		c.JSON(400, gin.H{"error": "用户已存在"})
		return err
	}
	return nil
}

// 注册缓存设置
func Register_write(c *gin.Context, u *model.User) error {
	if err := box.Temp.RE.HSet(c, u.Username, "username", u.Username).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return err
	}
	box.Temp.RE.Expire(c, u.Username, 2*time.Hour)
	return nil
}

// 登录缓存检查与处理
func Login_read(c *gin.Context, u *model.User) error {
	userTemp, err := box.Temp.RE.HMGet(c, u.Username, "username", "password").Result()
	if err != nil {
		return err
	}
	if userTemp[0] == nil || userTemp[1] == nil {
		return errors.New("登录缓存，数据解析为空")
	}
	userName, ok1 := userTemp[0].(string)
	passWord, ok2 := userTemp[1].(string)
	if !ok1 || !ok2 {
		return errors.New("登录缓存，数据转换错误")
	}
	if err := box.Temp.CheckUser(userName, passWord); err != nil {
		return err
	}
	if token, errToken := utils.GenerateToken(userName); errToken != nil {
		return errToken
	} else {
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		return nil
	}

}

// 登录缓存设置与处理
func Login_write(c *gin.Context, u *model.User) error {
	if err := box.Temp.RE.HMSet(c, u.Username,
		"username", u.Username,
		"password", u.Password,
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return err
	}
	box.Temp.RE.Expire(c, u.Username, 2*time.Hour)
	return nil
}

// 添加头像缓存处理
func AddUserHand_write(c *gin.Context, user *model.User) error {
	if err := box.Temp.RE.HSet(c, user.Username, "userhand", user.UserHand).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return err
	}
	return nil
}

package service

import (
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 用户注册服务
func UserRegister(c *gin.Context, u *model.User) {
	if !redis.Register_read(c, u) {
		return
	}
	switch mysql.Register(c, u) {
	case 1:
		{
			redis.Register_write(c, u)
		}
	default:
		return
	}
}

// 用户登录服务
func UserLogin(c *gin.Context, u *model.User) {
	switch redis.Login_read(c, u) {
	case 0:
		{

			switch mysql.Login(c, u) {
			case 1:
				redis.Login_write(c, u)
			default:
				return
			}
		}
	default:
		return
	}
}

// 添加用户头像服务
func UserHand(c *gin.Context) {
	mysql.AddUserHand(c)
}

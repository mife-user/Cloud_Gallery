package api

import (
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(c *gin.Context) {
	var u model.User
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}

	if !redis.Register_read(c, &u) {
		return
	}
	switch mysql.Register(c, &u) {
	case 1:
		{
			redis.Register_write(c, &u)
		}
	default:
		return
	}

}

// 登录
func Login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	switch redis.Login_read(c, &u) {
	case 0:
		{

			switch mysql.Login(c, &u) {
			case 1:
				redis.Login_write(c, &u)
			default:
				return
			}
		}
	default:
		return
	}
}

// 添加头像
func AddUserHand(c *gin.Context) {
	mysql.AddUserHand(c)
}

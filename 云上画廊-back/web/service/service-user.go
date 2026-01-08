package service

import (
	"fmt"
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 用户注册服务
func UserRegister(c *gin.Context, u *model.User) {

	/*查看是否已存在*/
	if err := redis.Register_read(c, u); err != nil {
		/*失败打印*/
		fmt.Println(err.Error())
		return
	}

	/*调用数据库，失败则返回*/
	if errRegis := mysql.Register(c, u); errRegis != nil {
		fmt.Println(errRegis.Error())
		return
	}

	/*成功则写入缓存*/
	if errRegisWrite := redis.Register_write(c, u); errRegisWrite != nil {
		fmt.Println(errRegisWrite.Error())
		return
	}

}

// 用户登录服务
func UserLogin(c *gin.Context, u *model.User) {

	if err := redis.Login_read(c, u); err != nil {

		/*失败打印*/
		fmt.Println(err.Error())

		/*调用数据库，失败则返回*/
		if err := mysql.Login(c, u); err != nil {
			fmt.Println(err.Error())
			return
		}

		/*成功则写入缓存*/
		if errLoginWrite := redis.Login_write(c, u); errLoginWrite != nil {
			fmt.Println(errLoginWrite.Error())
			return
		}

	}

}

// 添加用户头像服务
func UserHand(c *gin.Context, user *model.User) {
	/*数据库添加头像，失败则返回*/
	if err := mysql.AddUserHand(c, user); err != nil {
		fmt.Println(err.Error())
		return
	}
	/*成功则写入数据库*/
	if err := redis.AddUserHand_write(c, user); err != nil {
		fmt.Println(err.Error())
		return
	}
}

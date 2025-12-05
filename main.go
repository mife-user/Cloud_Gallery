package main

import (
	"painting/api"
	"painting/dao"
	"painting/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 仓库管理员先上班，准备好账本
	dao.Init()

	// 2. 启动 Gin 引擎（默认带了 Logger 和 Recovery 中间件，防崩坏）
	r := gin.Default()
	// 告诉 Gin，静态文件就在当前目录下
	r.LoadHTMLFiles("云上画廊.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "云上画廊.html", nil)
	})

	// 3. 全局中间件：安排外交官（CORS）站在大门口，所有请求都要经过它
	r.Use(middleware.Cors())

	//公共区（不需要证件）
	r.POST("/register", api.Register) //注册
	r.POST("/login", api.Login)       //登录
	r.GET("/gallery/:who", api.View)  //看展
	// :who 是动态参数，匹配如 /gallery/jack

	// VIP区（必须有证件【密钥】）
	authGroup := r.Group("/my")
	/*！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！*/
	authGroup.Use(middleware.AuthMiddleware()) // 安排保安
	{
		authGroup.POST("/upload", api.UploadPaint)                     //发布作品
		authGroup.POST("/delect", api.Delect)                          //删除作品
		authGroup.POST("/comment", api.PostComment)                    //评论作品
		authGroup.POST("/delectothercomment", api.DelectCommentMaster) //删除他人评论
		authGroup.POST("/delectmycomment", api.DelectCommentPoster)    //删除自己评论
	}
	/*！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！*/
	// 4. 开业！默认监听 0.0.0.0:8080
	/*！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！*/
	r.Run(":8080")
	/*！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！*/
}

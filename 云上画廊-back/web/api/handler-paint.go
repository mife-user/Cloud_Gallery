package api

import (
	"net/http"
	"painting/model"
	"painting/web/service"

	"github.com/gin-gonic/gin"
)

// 挂画
func UploadPaint(c *gin.Context) {
	var work model.Work
	usernameI, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "未认证"})
		return
	}
	username, _ := usernameI.(string)
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title 必填"})
		return
	}
	work.Content = c.PostForm("content")
	work.Author = username
	work.Title = title
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "必须上传图片文件"})
		return
	}
	service.PaintUp(c, &work, file)
}

// 删画
func Delect(c *gin.Context) {
	var work model.Work
	username, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "身份验证失败"})
		return
	}
	name := username.(string)
	if err := c.ShouldBindJSON(&work); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	service.PaintDelect(c, &work, name)
}

// 看展
func View(c *gin.Context) {
	who := c.Param("who")
	service.PaintView(c, who)
}

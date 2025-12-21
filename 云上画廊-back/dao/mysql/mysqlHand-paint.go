package mysql

import (
	"fmt"
	"net/http"
	"os"
	"painting/box"
	"painting/dao/redis"
	"painting/model"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 数据库挂画
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

	ext := filepath.Ext(file.Filename)
	uniqueName := fmt.Sprintf("%s_%d%s", username, time.Now().UnixNano(), ext)
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return
	}
	filePath := filepath.Join("uploads", uniqueName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return
	}

	// 5. 写入 DB
	work.Image = "/uploads/" + uniqueName
	if ok := box.Temp.AddWork(username, &work); !ok {
		c.JSON(500, gin.H{"error": "添加作品到数据库失败"})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
	redis.UpPaint_write(c, &work)
}

// 数据库删除画
func DelectPaint(c *gin.Context, work *model.Work) {
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
	if work.Author != "" && work.Author != name {
		c.JSON(403, gin.H{"error": "没有权限删除别人的作品"})
		return
	}
	if box.Temp.DelectPaint(name, work.Title) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，找不到画"})
	}
	redis.DePaint_write(c, work)
}

// 数据库看画
func View(c *gin.Context, who string) {
	works, err := box.Temp.GetWorks(who)
	if err != nil {
		c.JSON(400, gin.H{"error": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"owner": who, "works": works})
	redis.View_write(c, &works)
}

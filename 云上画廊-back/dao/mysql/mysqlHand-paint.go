package mysql

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"painting/box"
	"painting/dao/redis"
	"painting/model"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 数据库挂画
func UploadPaint(c *gin.Context, work *model.Work, file *multipart.FileHeader) error {
	ext := filepath.Ext(file.Filename)
	uniqueName := fmt.Sprintf("%s_%d%s", work.Author, time.Now().UnixNano(), ext)
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return err
	}
	filePath := filepath.Join("uploads", uniqueName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return err
	}

	// 5. 写入 DB
	work.Image = "/uploads/" + uniqueName
	if err := box.Temp.AddWork(work.Author, work); err != nil {
		c.JSON(500, gin.H{"error": "添加作品到数据库失败"})
		return err
	}
	c.JSON(200, gin.H{"message": "ok"})
	return nil
}

// 数据库删除画
func DelectPaint(c *gin.Context, work *model.Work, name string) error {
	if work.Author != "" && work.Author != name {
		c.JSON(403, gin.H{"error": "没有权限删除别人的作品"})
		return errors.New("数据库删除画作权限错误")
	}
	if err := box.Temp.DelectPaint(name, work.Title); err != nil {
		c.JSON(400, gin.H{"error": "数据库删除失败"})
		return err
	}
	c.JSON(200, gin.H{"message": "删除成功"})
	return nil
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

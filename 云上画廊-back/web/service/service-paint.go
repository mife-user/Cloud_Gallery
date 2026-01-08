package service

import (
	"fmt"
	"mime/multipart"
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 挂画服务
func PaintUp(c *gin.Context, work *model.Work, file *multipart.FileHeader) {
	if err := mysql.UploadPaint(c, work, file); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := redis.UpPaint_write(c, work); err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 删画服务
func PaintDelect(c *gin.Context, w *model.Work, name string) {
	if err := mysql.DelectPaint(c, w, name); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := redis.DePaint_write(c, w); err != nil {
		fmt.Println(err.Error())
		return
	}
}

// 看展服务
func PaintView(c *gin.Context, who string) {
	if redis.View_read(c, who) {
		return
	}
	mysql.View(c, who)
}

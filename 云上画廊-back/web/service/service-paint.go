package service

import (
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 挂画服务
func PaintUp(c *gin.Context) {
	mysql.UploadPaint(c)
}

// 删画服务
func PaintDelect(c *gin.Context, w *model.Work) {
	mysql.DelectPaint(c, w)
}

// 看展服务
func PaintView(c *gin.Context, who string) {
	if redis.View_read(c, who) {
		return
	}
	mysql.View(c, who)
}

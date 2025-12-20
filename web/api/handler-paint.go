package api

import (
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 挂画
func UploadPaint(c *gin.Context) {
	mysql.UploadPaint(c)
}

// 删画
func Delect(c *gin.Context) {
	var w model.Work
	mysql.DelectPaint(c, &w)
}

// 看展
func View(c *gin.Context) {
	who := c.Param("who")
	if redis.View_read(c, who) {
		return
	}
	mysql.View(c, who)
}

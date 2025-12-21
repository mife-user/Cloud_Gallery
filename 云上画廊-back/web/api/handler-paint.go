package api

import (
	"painting/model"
	"painting/web/service"

	"github.com/gin-gonic/gin"
)

// 挂画
func UploadPaint(c *gin.Context) {
	service.PaintUp(c)
}

// 删画
func Delect(c *gin.Context) {
	var w model.Work
	service.PaintDelect(c, &w)
}

// 看展
func View(c *gin.Context) {
	who := c.Param("who")
	service.PaintView(c, who)
}

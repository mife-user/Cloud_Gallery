package api

import (
	"painting/model"
	"painting/web/service"

	"github.com/gin-gonic/gin"
)

// 评论
func PostComment(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "登录名错误"})
		return
	}
	commentator := username.(string)
	var req model.CommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数没填对"})
		return
	}
	newComment := model.Comment{
		FromUser: commentator,
		Content:  req.Content,
	}
	service.CommentPost(c, &newComment, &req)

}

// 作者删除评论
func DelectCommentMaster(c *gin.Context) {
	/*---------------------------------------------------------*/
	user, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "当前账号出现问题"})
		return
	}
	currentMaster := user.(string)
	var req model.DeleteCommentReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数解析失败，请检查数据格式"})
		return
	}
	if req.Owner != currentMaster {
		c.JSON(403, gin.H{"error": "只能删除自己作品下的评论", "message": currentMaster + "!=" + req.Owner})
		return
	}
	comment := model.Comment{
		FromUser: req.FromUser,
		Content:  req.Content,
	}
	/*---------------------------------------------------------*/
	service.DelectCommentMaster(c, currentMaster, &req, &comment)
}

// 用户删除评论
func DelectCommentPoster(c *gin.Context) {
	user, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "用户未登录"})
		return
	}
	currentUser := user.(string)
	var req model.DeleteCommentReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数格式错误"})
		return
	}
	if req.FromUser != currentUser {
		c.JSON(403, gin.H{"error": "你没有权限删除别人的评论"})
		return
	}
	comment := model.Comment{
		FromUser: req.FromUser,
		Content:  req.Content,
	}
	service.DelectCommentPoster(c, &req, &comment)
}

package mysql

import (
	"painting/box"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 数据库评论
func PostComment(c *gin.Context, newComment *model.Comment, req *model.CommentRequest) bool {
	if box.Temp.AddComment(req.TargetAuthor, req.WorkTitle, newComment) {
		c.JSON(200, gin.H{"message": "评论成功", "data": newComment})
		return true
	} else {
		c.JSON(404, gin.H{"error": "找不到这幅画"})
	}
	return false
}

// 数据库作者删除评论
func DelectCommentMaster(c *gin.Context, currentMaster string, workTitle string, comment *model.Comment) {
	if box.Temp.DelectComment(currentMaster, workTitle, comment) {
		c.JSON(200, gin.H{"message": "作为作者，已删除该评论"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，未找到该评论或画作"})
	}
}

// 数据库用户删除评论
func DelectCommentPoster(c *gin.Context, req *model.DeleteCommentReq, comment *model.Comment) {
	if box.Temp.DelectComment(req.Owner, req.Title, comment) {
		c.JSON(200, gin.H{"message": "已撤回您的评论"})
	} else {
		c.JSON(400, gin.H{"error": "撤回失败，可能评论已不存在"})
	}
}

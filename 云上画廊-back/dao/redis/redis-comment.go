package redis

import (
	"fmt"
	"painting/box"
	"painting/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 评论缓存处理
func PostComment_read(c *gin.Context, newComment *model.Comment, req *model.CommentRequest) bool {
	workComment := fmt.Sprintf("%s:%s:%s:%s:comments",
		req.TargetAuthor, req.WorkTitle, newComment.FromUser, newComment.CreatedAt.Format(time.RFC3339))

	if err := box.Temp.RE.HMSet(c,
		workComment,
		"id", strconv.FormatUint(uint64(newComment.ID), 10),
		"from_user", newComment.FromUser,
		"content", newComment.Content,
		"created_at", newComment.CreatedAt.Format(time.RFC3339),
		"updated_at", newComment.UpdatedAt.Format(time.RFC3339),
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return false
	}
	box.Temp.RE.Expire(c, workComment, 2*time.Hour)
	return true
}

// 作者评论删除缓存处理
func DelectCommentMaster_read(c *gin.Context, currentMaster string, req *model.DeleteCommentReq) {

	workComment := fmt.Sprintf("%s:%s:%s:%s:comments", currentMaster, req.Title, req.FromUser, req.CreatedAt.Format(time.RFC3339))
	if err := box.Temp.RE.HDel(c,
		workComment,
		"from_user",
		"content",
		"created_at",
		"updated_at",
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

// 用户评论删除缓存处理
func DelectCommentPoster_read(c *gin.Context, req *model.DeleteCommentReq) {
	workComment := fmt.Sprintf("%s:%s:%s:%s:comments", req.Owner, req.Title, req.FromUser, req.CreatedAt.Format(time.RFC3339))
	if err := box.Temp.RE.HDel(c,
		workComment,
		"from_user",
		"content",
		"created_at",
		"updated_at",
	).Err(); err != nil {
		c.JSON(400, gin.H{"error": "redis服务器错误"})
		return
	}
}

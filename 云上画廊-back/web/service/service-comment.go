package service

import (
	"painting/dao/mysql"
	"painting/dao/redis"
	"painting/model"

	"github.com/gin-gonic/gin"
)

// 发布评论服务
func CommentPost(c *gin.Context, newComment *model.Comment, req *model.CommentRequest) {
	mysql.PostComment(c, newComment, req)
	redis.PostComment_read(c, newComment, req)

}

// 作者删除评论服务
func DelectCommentMaster(c *gin.Context, currentMaster string, req *model.DeleteCommentReq, comment *model.Comment) {
	mysql.DelectCommentMaster(c, currentMaster, req.Title, comment)
	redis.DelectCommentMaster_read(c, currentMaster, req)
}

// 用户删除评论服务
func DelectCommentPoster(c *gin.Context, req *model.DeleteCommentReq, comment *model.Comment) {
	mysql.DelectCommentPoster(c, req, comment)
	redis.DelectCommentPoster_read(c, req)
}

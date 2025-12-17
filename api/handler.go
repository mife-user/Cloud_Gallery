package api

import (
	"fmt"
	"painting/dao"
	"painting/model"
	"time"

	"github.com/gin-gonic/gin"
)

var Temp *dao.Database

// 关闭数据库
func CloseSQL() {
	if err := Temp.Close(); err != nil {
		fmt.Print(err)
		return
	}

}

// 注册
func Register(c *gin.Context) {
	var u model.User
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}

	if !register1(c, &u) {
		return
	}
	switch register(c, &u) {
	case 1:
		{
			register2(c, &u)
		}
	default:
		return
	}

}

// 登录
func Login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBind(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	switch login1(c, &u) {
	case 0:
		{

			switch login(c, &u) {
			case 1:
				login2(c, &u)
			default:
				return
			}
		}
	default:
		return
	}
}

// 挂画
func UploadPaint(c *gin.Context) {
	uploadPaint(c)
}

// 删画
func Delect(c *gin.Context) {
	var w model.Work
	delectPaint(c, &w)
}

// 看展
func View(c *gin.Context) {
	who := c.Param("who")
	if view1(c, who) {
		return
	}
	view(c, who)
}

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
		FromUser:  commentator,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if postComment(c, &newComment, &req) {
		postComment1(c, &newComment, &req)
	}

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
		c.JSON(403, gin.H{"error": "只能删除自己作品下的评论"})
		return
	}
	comment := model.Comment{
		FromUser:  req.FromUser,
		Content:   req.Content,
		CreatedAt: req.CreatedAt,
	}
	/*---------------------------------------------------------*/
	delectCommentMaster1(c, currentMaster, &req)
	delectCommentMaster(c, currentMaster, req.Title, &comment)
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
		FromUser:  req.FromUser,
		Content:   req.Content,
		CreatedAt: req.CreatedAt,
	}
	delectCommentPoster1(c, &req)
	delectCommentPoster(c, &req, &comment)
}

// 添加头像
func AddUserHand(c *gin.Context) {
	addUserHand(c)
}

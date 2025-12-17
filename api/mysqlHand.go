package api

import (
	"fmt"
	"net/http"
	"os"
	"painting/model"
	"painting/utils"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 数据库注册
func register(c *gin.Context, u *model.User) int {
	if ok := Temp.AddUser(u.Username, u.Password); ok {
		c.JSON(200, gin.H{"message": "注册成功！"})
		return 1
	}
	return 2
}

// 数据库登录
func login(c *gin.Context, u *model.User) int {
	if ok := Temp.CheckUser(u.Username, u.Password); ok {
		token, _ := utils.GenerateToken(u.Username)
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
		return 1
	} else {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return 2
	}
}

// 数据库挂画
func uploadPaint(c *gin.Context) {
	var work model.Work
	usernameI, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "未认证"})
		return
	}
	username, _ := usernameI.(string)
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title 必填"})
		return
	}
	work.Content = c.PostForm("content")
	work.Author = username
	work.Title = title
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "必须上传图片文件"})
		return
	}

	ext := filepath.Ext(file.Filename)
	uniqueName := fmt.Sprintf("%s_%d%s", username, time.Now().UnixNano(), ext)
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return
	}
	filePath := filepath.Join("uploads", uniqueName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return
	}

	// 5. 写入 DB
	work.Image = "/uploads/" + uniqueName
	if ok := Temp.AddWork(username, &work); !ok {
		c.JSON(500, gin.H{"error": "添加作品到数据库失败"})
		return
	}
	c.JSON(200, gin.H{"message": "ok"})
	uploadPaint2(c, &work)
}

// 数据库删除画
func delectPaint(c *gin.Context, work *model.Work) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "身份验证失败"})
		return
	}
	name := username.(string)
	if err := c.ShouldBindJSON(&work); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	if work.Author != "" && work.Author != name {
		c.JSON(403, gin.H{"error": "没有权限删除别人的作品"})
		return
	}
	if Temp.DelectPaint(name, work.Title) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，找不到画"})
	}
	delectPaint2(c, work)
}

// 数据库看画
func view(c *gin.Context, who string) {
	works, err := Temp.GetWorks(who)
	if err != nil {
		c.JSON(400, gin.H{"error": "数据库错误"})
		return
	}
	c.JSON(200, gin.H{"owner": who, "works": works})
	view2(c, &works)
}

// 数据库评论
func postComment(c *gin.Context, newComment *model.Comment, req *model.CommentRequest) bool {
	if Temp.AddComment(req.TargetAuthor, req.WorkTitle, newComment) {
		c.JSON(200, gin.H{"message": "评论成功", "data": newComment})
		return true
	} else {
		c.JSON(404, gin.H{"error": "找不到这幅画"})
	}
	return false
}

// 数据库作者删除评论
func delectCommentMaster(c *gin.Context, currentMaster string, workTitle string, comment *model.Comment) {
	if Temp.DelectComment(currentMaster, workTitle, comment) {
		c.JSON(200, gin.H{"message": "作为作者，已删除该评论"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，未找到该评论或画作"})
	}
}

// 数据库用户删除评论
func delectCommentPoster(c *gin.Context, req *model.DeleteCommentReq, comment *model.Comment) {
	if Temp.DelectComment(req.Owner, req.Title, comment) {
		c.JSON(200, gin.H{"message": "已撤回您的评论"})
	} else {
		c.JSON(400, gin.H{"error": "撤回失败，可能评论已不存在"})
	}
}

// 数据库添加头像
func addUserHand(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, gin.H{"error": "未登录"})
		return
	}
	if user.Username == "" {
		c.JSON(400, gin.H{"error": "username required"})
		return
	}
	file, err := c.FormFile("userhand")
	if err != nil {
		c.JSON(400, gin.H{"error": "file required"})
		return
	}
	dstName := fmt.Sprintf("%s_%d_%s", user.Username, time.Now().Unix(), filepath.Ext(file.Filename))
	if err := os.MkdirAll("userhands", os.ModePerm); err != nil {
		c.JSON(500, gin.H{"error": "创建目录失败"})
		return
	}
	filePath := filepath.Join("userhands", dstName)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return
	}
	c.JSON(200, gin.H{"message": "ok", "path": "/userhands/" + dstName})
	user.UserHand = "/userhands/" + dstName
	Temp.AddHand(&user)
	addUserHand2(c, &user)
}

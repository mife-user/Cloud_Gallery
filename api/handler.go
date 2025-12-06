package api

import (
	"fmt"
	"os"
	"painting/dao"
	"painting/model"
	"painting/utils"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// 注册
func Register(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	if dao.AddUser(u.Username, u.Password) {
		c.JSON(200, gin.H{"message": "注册成功！"})
	} else {
		c.JSON(400, gin.H{"error": "用户名已存在"})
	}
}

// 登录
func Login(c *gin.Context) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	if dao.CheckUser(u.Username, u.Password) {
		token, _ := utils.GenerateToken(u.Username)
		c.JSON(200, gin.H{
			"message": "登录成功",
			"token":   token,
		})
	} else {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
	}
}

// 挂画
// 在UploadPaint函数中修改图片保存逻辑
func UploadPaint(c *gin.Context) {
	username := c.GetString("username")
	title := c.PostForm("title")
	content := c.PostForm("content")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "必须上传图片文件"})
		return
	}

	// 创建唯一文件名避免冲突
	ext := filepath.Ext(file.Filename)
	uniqueName := fmt.Sprintf("%s_%d%s", username, time.Now().UnixNano(), ext)

	// 保存到本地 uploads 目录
	os.MkdirAll("uploads", 0777)
	filePath := "uploads/" + uniqueName

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{"error": "保存文件失败"})
		return
	}

	// 存储相对路径，前端通过相对路径访问
	// 如果是远程服务器，需要配置静态文件服务
	work := model.Work{
		Title:   title,
		Image:   "/uploads/" + uniqueName, // 使用相对URL路径
		Content: content,
	}

	dao.AddWork(username, work)
	c.JSON(200, gin.H{"msg": "ok"})
}

// 删画
func Delect(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "身份验证失败"})
		return
	}
	// 修正：这里之前的 err1 命名和逻辑有点乱，整理了一下
	name := username.(string)

	var w model.Work
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}
	// 只能删自己的画
	if dao.DelectPaint(name, w.Title) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，找不到画"})
	}
}

// 看展
func View(c *gin.Context) {
	who := c.Param("who")
	works := dao.GetWorks(who)
	c.JSON(200, gin.H{"owner": who, "works": works})
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数没填对"})
		return
	}
	newComment := model.Comment{
		FromUser:  commentator,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if dao.AddComment(req.TargetAuthor, req.WorkTitle, newComment) {
		c.JSON(200, gin.H{"message": "评论成功", "data": newComment})
	} else {
		c.JSON(404, gin.H{"error": "找不到这幅画"})
	}
}

// 作者删除评论 (我是画的主人，我看不惯这条评论)
func DelectCommentMaster(c *gin.Context) {
	user, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "当前账号出现问题"})
		return
	}
	currentMaster := user.(string) // 当前登录的人

	// [修改] 使用统一的 DeleteCommentReq 接收所有参数
	var req model.DeleteCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数解析失败，请检查数据格式"})
		return
	}

	// [逻辑校验] 只有画廊的主人自己才能行使 "作者删除权"
	// 前端虽然没传 owner，但默认当前登录者就是 owner，或者我们可以校验 req.Owner
	// 这里我们直接认为：你要删你名下的画的评论，那你必须是 currentMaster

	// 组装评论对象用于查找
	comment := model.Comment{
		FromUser:  req.FromUser,
		Content:   req.Content,
		CreatedAt: req.CreatedAt,
	}

	// 这里的第一个参数传 currentMaster，确保是在操作自己的画廊
	if dao.DelectComment(currentMaster, req.Title, comment) {
		c.JSON(200, gin.H{"message": "作为作者，已删除该评论"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，未找到该评论或画作"})
	}
}

// 用户删除评论 (我自己写的评论，我想撤回)
func DelectCommentPoster(c *gin.Context) {
	user, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "用户未登录"})
		return
	}
	currentUser := user.(string) // 当前登录的人

	// [修改] 之前试图分两次 Bind 是错误的，必须一次性 Bind
	// 同时也去掉了对 c.Get("commenttime") 的依赖，因为中间件里根本没存这个
	var req model.DeleteCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "参数格式错误"})
		return
	}

	// [逻辑校验] 只有评论的作者(FromUser) 和 当前登录人(currentUser) 一致，才有资格撤回
	if req.FromUser != currentUser {
		c.JSON(403, gin.H{"error": "你没有权限删除别人的评论"})
		return
	}

	comment := model.Comment{
		FromUser:  req.FromUser,
		Content:   req.Content,
		CreatedAt: req.CreatedAt,
	}

	// 注意：这里的第一个参数是 req.Owner (画挂在谁家)，而不是 currentUser
	// 因为我们要去 req.Owner 的画廊里，找到这幅画，删掉 currentUser 写的评论
	if dao.DelectComment(req.Owner, req.Title, comment) {
		c.JSON(200, gin.H{"message": "已撤回您的评论"})
	} else {
		c.JSON(400, gin.H{"error": "撤回失败，可能评论已不存在"})
	}
}

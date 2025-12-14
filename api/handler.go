package api

import (
	"fmt"
	"net/http"
	"os"
	"painting/dao"
	"painting/model"
	"painting/utils"
	"path/filepath"
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
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	if ok := Temp.AddUser(u.Username, u.Password); ok {
		c.JSON(200, gin.H{"message": "注册成功！"})
	} else {
		c.JSON(400, gin.H{"error": "用户名已存在"})
	}
}

// 登录
func Login(c *gin.Context) {
	var u model.User
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{"error": "参数填错了"})
		return
	}
	if ok := Temp.CheckUser(u.Username, u.Password); ok {
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
	usernameI, ok := c.Get("username")
	if !ok {
		c.JSON(401, gin.H{"error": "未认证"})
		return
	}
	username, _ := usernameI.(string)
	var work model.Work
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title 必填"})
		return
	}
	work.Content = c.PostForm("content")
	work.Author = c.PostForm("author")
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

}

// 删画
func Delect(c *gin.Context) {
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
	username, ok := c.Get("username")
	if !ok {
		c.JSON(400, gin.H{"error": "身份验证失败"})
		return
	}
	name := username.(string)

	var w model.Work
	if err := c.ShouldBindJSON(&w); err != nil {
		c.JSON(400, gin.H{"error": "参数错误"})
		return
	}

	// 双重验证：确保只能删自己的画
	// 这里前端传过来的 w.Author 应该是当前用户的用户名
	if w.Author != "" && w.Author != name {
		c.JSON(403, gin.H{"error": "没有权限删除别人的作品"})
		return
	}

	// 只能删自己的画
	if Temp.DelectPaint(name, w.Title) {
		c.JSON(200, gin.H{"message": "删除成功"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，找不到画"})
	}
}

// 看展
func View(c *gin.Context) {
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
	who := c.Param("who")
	works := Temp.GetWorks(who)
	c.JSON(200, gin.H{"owner": who, "works": works})
}

// 评论
func PostComment(c *gin.Context) {
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
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
	if Temp.AddComment(req.TargetAuthor, req.WorkTitle, newComment) {
		c.JSON(200, gin.H{"message": "评论成功", "data": newComment})
	} else {
		c.JSON(404, gin.H{"error": "找不到这幅画"})
	}
}

// 作者删除评论 (我是画的主人，我看不惯这条评论)
func DelectCommentMaster(c *gin.Context) {
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
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
	comment := model.Comment{
		FromUser:  req.FromUser,
		Content:   req.Content,
		CreatedAt: req.CreatedAt,
	}
	if Temp.DelectComment(currentMaster, req.Title, comment) {
		c.JSON(200, gin.H{"message": "作为作者，已删除该评论"})
	} else {
		c.JSON(400, gin.H{"error": "删除失败，未找到该评论或画作"})
	}
}

// 用户删除评论 (我自己写的评论，我想撤回)
func DelectCommentPoster(c *gin.Context) {
	// temp, ok := dao.Init()
	// if !ok {
	// 	c.JSON(400, gin.H{"error": "数据库错误"})
	// 	return
	// }
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
	if Temp.DelectComment(req.Owner, req.Title, comment) {
		c.JSON(200, gin.H{"message": "已撤回您的评论"})
	} else {
		c.JSON(400, gin.H{"error": "撤回失败，可能评论已不存在"})
	}
}

// 添加头像
func AddUserHand(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
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
	// 确保目录存在并保存文件，生成唯一名
	dstName := fmt.Sprintf("%s_%d_%s", user.Username, time.Now().Unix(), filepath.Base(file.Filename))
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
}
